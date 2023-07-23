package chatroom

import (
	"context"
	"sync"
)

type MessageChannel chan *Message

type EventChannel chan *ParticipantEvent

func NewMessageChannel() MessageChannel {
	return make(chan *Message)
}

func NewEventChannel() EventChannel {
	return make(chan *ParticipantEvent)
}

type MessageListenerFunc func(*Message)

type EventListenerFunc func(*ParticipantEvent)

type Participant struct {
	id             string
	name           string
	messageChannel MessageChannel
	eventChannel   EventChannel
	room           *Room
	alive          bool
	mu             sync.Mutex
}

func (p *Participant) Room() *Room {
	return p.room
}

func (p *Participant) SetRoom(room *Room) {
	p.room = room
}

func NewParticipant(id string, name string) *Participant {
	return &Participant{
		id:             id,
		name:           name,
		messageChannel: NewMessageChannel(),
		eventChannel:   NewEventChannel(),
		alive:          true,
		mu:             sync.Mutex{},
	}
}

func (p *Participant) Id() string {
	return p.id
}

func (p *Participant) Name() string {
	return p.name
}

func (p *Participant) OnListen(ctx context.Context, listener MessageListenerFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-p.messageChannel:
			listener(m)
		}
	}
}

func (p *Participant) OnEvent(ctx context.Context, listener EventListenerFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-p.eventChannel:
			listener(e)
		}
	}
}

func (p *Participant) Leave() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.alive {
		return
	}
	p.room.RemoveParticipant(p)
	close(p.messageChannel)
	close(p.eventChannel)
	p.alive = false
}

func (p *Participant) Say(m *Message) {
	m.setParticipant(p)
	p.room.Broadcast(m)
}
