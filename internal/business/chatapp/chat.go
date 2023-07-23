package chatapp

import (
	"chatty/internal/business/chatroom"
	"chatty/internal/platform/id"
	"errors"
	"fmt"
)

var ErrCannotJoinRoomNotFound = errors.New("cannot join room: room not found")
var ErrCannotJoinRoomInternal = errors.New("cannot join room: internal error")
var ErrCannotLeaveRoomParticipantNotAssigned = errors.New("cannot leave room: participant not assigned")
var ErrCannotLeaveRoomInternal = errors.New("cannot leave room: internal error")

type App struct {
	roomCollection chatroom.Collection
	idGenerator    id.Generator
}

func NewChatApp(roomCollection chatroom.Collection, idGenerator id.Generator) *App {
	return &App{roomCollection: roomCollection, idGenerator: idGenerator}
}

func (c App) CreateRoom(name string) (*chatroom.Room, error) {
	roomId := c.idGenerator.MustGenerate()
	newRoom := chatroom.NewRoom(roomId, name)

	err := c.roomCollection.Save(newRoom)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (c App) JoinRoom(roomId string, participantName string) (*chatroom.Participant, error) {
	room, err := c.roomCollection.GetByID(roomId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCannotJoinRoomNotFound, err)
	}

	participantId := c.idGenerator.MustGenerate()
	participant := chatroom.NewParticipant(participantId, participantName)
	room.AddParticipant(participant)
	err = c.roomCollection.Save(room)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCannotJoinRoomInternal, err)
	}
	return participant, nil
}

func (c App) LeaveRoom(participant *chatroom.Participant) error {
	room := participant.Room()
	if room != nil {
		return ErrCannotLeaveRoomParticipantNotAssigned
	}
	room.RemoveParticipant(participant)
	err := c.roomCollection.Save(room)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCannotLeaveRoomInternal, err)
	}
	return nil
}
