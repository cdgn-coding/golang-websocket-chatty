package chatroom

type Message struct {
	text        string
	participant *Participant
}

func (m *Message) Participant() *Participant {
	return m.participant
}

func (m *Message) setParticipant(participant *Participant) {
	m.participant = participant
}

func (m Message) Text() string {
	return m.text
}

type MessageOption func(*Message)

func WithParticipant(p *Participant) MessageOption {
	return func(m *Message) {
		m.participant = p
	}
}

func NewMessage(text string, options ...MessageOption) *Message {
	m := &Message{text: text}
	for _, option := range options {
		option(m)
	}
	return m
}
