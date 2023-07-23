package websocket

import (
	"chatty/internal/business/chatroom"
	"errors"
)

type NewMessage struct {
	Text        string `json:"text"`
	Participant string `json:"participant"`
}

func NewNewMessage(text string, participant string) *NewMessage {
	return &NewMessage{Text: text, Participant: participant}
}

const (
	ParticipantJoined string = "participant_joined"
	ParticipantLeft   string = "participant_left"
)

type ParticipantEvent struct {
	EventType   string `json:"event_type"`
	Participant string `json:"participant"`
}

func NewParticipantEvent(eventType string, participant string) *ParticipantEvent {
	return &ParticipantEvent{EventType: eventType, Participant: participant}
}

var ErrInvalidEventType = errors.New("invalid event type")

func MapEvent(event *chatroom.ParticipantEvent) (*ParticipantEvent, error) {
	switch event.EventType() {
	case chatroom.ParticipantJoined:
		return NewParticipantEvent(ParticipantJoined, event.Participant().Name()), nil
	case chatroom.ParticipantLeft:
		return NewParticipantEvent(ParticipantLeft, event.Participant().Name()), nil
	default:
		return nil, ErrInvalidEventType
	}
}

type ServerEventType string

const (
	ServerShuttingDown ServerEventType = "server_shutting_down"
)

type ServerEvent struct {
	EventType ServerEventType `json:"event_type"`
}

func NewServerEvent(eventType ServerEventType) *ServerEvent {
	return &ServerEvent{EventType: eventType}
}

func NewServerShuttingDown() *ServerEvent {
	return NewServerEvent(ServerShuttingDown)
}
