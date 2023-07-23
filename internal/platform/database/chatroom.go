package database

import (
	"chatty/internal/business/chatroom"
	"errors"
)

var ErrRoomNotFound = errors.New("room not found")

type InMemoryChatRoomCollection struct {
	rooms map[string]*chatroom.Room
}

func NewInMemoryChatRoomCollection() *InMemoryChatRoomCollection {
	return &InMemoryChatRoomCollection{rooms: make(map[string]*chatroom.Room)}
}

func (c InMemoryChatRoomCollection) GetByID(id string) (*chatroom.Room, error) {
	room, ok := c.rooms[id]
	if !ok {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (c InMemoryChatRoomCollection) Save(room *chatroom.Room) error {
	c.rooms[room.Id()] = room
	return nil
}

func (c InMemoryChatRoomCollection) Delete(room *chatroom.Room) error {
	delete(c.rooms, room.Id())
	return nil
}
