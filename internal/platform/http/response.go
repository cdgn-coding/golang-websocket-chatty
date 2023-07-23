package http

import "chatty/internal/business/chatroom"

type RoomCreated struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func MapResponseRoomCreated(room *chatroom.Room) *RoomCreated {
	return &RoomCreated{
		ID:   room.Id(),
		Name: room.Name(),
	}
}

func NewRoomCreated(ID string) *RoomCreated {
	return &RoomCreated{ID: ID}
}
