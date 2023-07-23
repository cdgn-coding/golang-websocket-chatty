package websocket

import (
	"chatty/internal/business/chatapp"
	"chatty/internal/business/chatroom"
	_ "chatty/internal/business/chatroom"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{}

type App struct {
	chatApp     *chatapp.App
	closing     bool
	connections uint32
	mu          sync.Mutex
}

func NewApp(app *chatapp.App) *App {
	return &App{
		chatApp: app,
		closing: false,
		mu:      sync.Mutex{},
	}
}

func (app *App) Shutdown(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	app.closing = true
	log.Printf("Closing all connections and processes...")
	for {
		select {
		case <-ctx.Done():
			log.Printf("Forcing shutdown websocket app - context done")
			return
		case <-ticker.C:
			if app.connections == 0 {
				log.Printf("Websocket app successfully shutdown")
				return
			}
		}
	}
}

func (app *App) Handle(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Cannot upgrade connection - %v", err)
		return
	}

	name := r.URL.Query().Get("name")
	room := r.URL.Query().Get("room")

	participant, err := app.chatApp.JoinRoom(room, name)
	if err != nil {
		log.Printf("Cannot join room - %v", err)
		return
	}

	app.OnConnect(connection)
	ctx, cancel := context.WithCancel(context.Background())
	defer app.handleDisconnect(cancel, participant, connection)
	go app.handleLiveness(cancel, connection)

	go participant.OnListen(ctx, app.handleNewMessage(connection))
	go participant.OnEvent(ctx, app.handleNewEvent(connection))

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("connection closed - %v", err)
			}
			break
		}

		switch mt {
		case websocket.TextMessage:
			chatMessage := chatroom.NewMessage(string(message))
			participant.Say(chatMessage)
		}
	}
}

func (app *App) handleLiveness(cancel context.CancelFunc, connection *websocket.Conn) {
	waitPeriod := 5 * time.Second

	connection.SetPongHandler(func(string) error {
		return connection.SetReadDeadline(time.Now().Add(waitPeriod))
	})

	go func() {
		ticker := time.NewTicker(waitPeriod)
		defer ticker.Stop()
		for {
			<-ticker.C
			if app.closing {
				cancel()
				event := NewServerShuttingDown()
				bytes, err := json.Marshal(event)
				if err != nil {
					log.Printf("Cannot marshal event - %v", err)
					return
				}

				waitPeriod := 3 * time.Second
				err = app.OnWrite(connection, bytes, waitPeriod)
				if err != nil {
					log.Printf("Cannot write message - %v", err)
					return
				}
				app.OnDisconnect(connection)
				return
			}
			err := connection.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("Error sending ping - %v", err)
				return
			}
		}
	}()
}

func (app *App) handleDisconnect(cancel context.CancelFunc, participant *chatroom.Participant, connection *websocket.Conn) {
	cancel()
	participant.Leave()
	app.OnDisconnect(connection)
}

func (app *App) handleNewEvent(connection *websocket.Conn) chatroom.EventListenerFunc {
	return func(e *chatroom.ParticipantEvent) {
		event, err := MapEvent(e)
		if err != nil {
			log.Printf("Cannot convert event - %v", err)
		}
		bytes, err := json.Marshal(event)
		if err != nil {
			log.Printf("Cannot marshal event - %v", err)
		}
		waitTime := 5 * time.Second
		app.OnWrite(connection, bytes, waitTime)
	}
}

func (app *App) handleNewMessage(connection *websocket.Conn) func(message *chatroom.Message) {
	return func(message *chatroom.Message) {
		event := NewNewMessage(message.Text(), message.Participant().Name())
		bytes, err := json.Marshal(event)
		if err != nil {
			log.Printf("Cannot marshal event - %v", err)
			return
		}
		waitTime := 5 * time.Second
		app.OnWrite(connection, bytes, waitTime)
		if err != nil {
			log.Printf("Cannot send message - %v", err)
		}
	}
}

func (app *App) OnWrite(connection *websocket.Conn, bytes []byte, deadline time.Duration) error {
	err := connection.SetWriteDeadline(time.Now().Add(deadline))
	if err != nil {
		return err
	}
	return connection.WriteMessage(websocket.TextMessage, bytes)
}

func (app *App) OnConnect(connection *websocket.Conn) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.connections++
}

func (app *App) OnDisconnect(connection *websocket.Conn) {
	app.mu.Lock()
	defer app.mu.Unlock()
	err := connection.Close()
	if err != nil {
		return
	}
	app.connections--
}
