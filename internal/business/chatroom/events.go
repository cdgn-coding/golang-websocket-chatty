package chatroom

type EventType string

const (
	ParticipantJoined EventType = "participant_joined"
	ParticipantLeft   EventType = "participant_left"
)

type ParticipantEvent struct {
	eventType   EventType
	participant *Participant
}

func NewParticipantJoined(participant *Participant) *ParticipantEvent {
	return &ParticipantEvent{eventType: ParticipantJoined, participant: participant}
}

func NewParticipantLeft(participant *Participant) *ParticipantEvent {
	return &ParticipantEvent{eventType: ParticipantLeft, participant: participant}
}

func (p ParticipantEvent) EventType() EventType {
	return p.eventType
}

func (p ParticipantEvent) Participant() *Participant {
	return p.participant
}
