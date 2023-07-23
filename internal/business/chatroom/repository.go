package chatroom

type Collection interface {
	GetByID(id string) (*Room, error)
	Save(room *Room) error
	Delete(room *Room) error
}
