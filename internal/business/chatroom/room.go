package chatroom

import "sync"

type Room struct {
	id           string
	name         string
	mu           sync.Mutex
	participants []*Participant
}

func NewRoom(id string, name string) *Room {
	return &Room{
		id:           id,
		name:         name,
		mu:           sync.Mutex{},
		participants: []*Participant{},
	}
}

func (r *Room) Id() string {
	return r.id
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Participants() []*Participant {
	return r.participants
}

func (r *Room) AddParticipant(p *Participant) {
	r.mu.Lock()
	p.room = r
	r.participants = append(r.participants, p)
	r.mu.Unlock()
	go r.BroadcastEvent(NewParticipantJoined(p))
}

func (r *Room) RemoveParticipant(p *Participant) {
	r.mu.Lock()
	for i, participant := range r.participants {
		if participant.id == p.id {
			r.participants = append(r.participants[:i], r.participants[i+1:]...)
			break
		}
	}
	r.mu.Unlock()
	r.BroadcastEvent(NewParticipantLeft(p))
}

func (r *Room) Broadcast(m *Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.participants {
		p.messageChannel <- m
	}
}

func (r *Room) BroadcastEvent(event *ParticipantEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.participants {
		p.eventChannel <- event
	}
}
