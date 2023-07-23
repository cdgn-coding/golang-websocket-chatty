package http

import (
	"chatty/internal/business/chatapp"
	"encoding/json"
	"log"
	"net/http"
)

type HttpApp struct {
	app *chatapp.App
}

func NewHttpApp(app *chatapp.App) *HttpApp {
	return &HttpApp{app: app}
}

func (a *HttpApp) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (a *HttpApp) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, err := a.app.CreateRoom(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := MapResponseRoomCreated(room)
	bytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Cannot marshal response - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
