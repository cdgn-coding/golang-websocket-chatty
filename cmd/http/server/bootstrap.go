package main

import (
	"chatty/internal/business/chatapp"
	"chatty/internal/platform/database"
	"chatty/internal/platform/http"
	"chatty/internal/platform/id"
	"chatty/internal/platform/websocket"
)

type Container struct {
	roomsDatabase *database.InMemoryChatRoomCollection
	idGenerator   *id.KsuidGenerator
	app           *chatapp.App
	websocketApp  *websocket.App
	httpApp       *http.HttpApp
}

func bootstrap() *Container {
	roomsDatabase := database.NewInMemoryChatRoomCollection()
	idGenerator := id.NewKsuidGenerator()
	app := chatapp.NewChatApp(roomsDatabase, idGenerator)
	websocketApp := websocket.NewApp(app)
	httpApp := http.NewHttpApp(app)
	return &Container{
		roomsDatabase: roomsDatabase,
		idGenerator:   idGenerator,
		app:           app,
		websocketApp:  websocketApp,
		httpApp:       httpApp,
	}
}
