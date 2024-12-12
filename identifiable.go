package ebui

type Identifiable interface {
	GetID() uint64
}

var nextID uint64

func GenerateID() uint64 {
	nextID++
	if nextID == 0 {
		panic("ID overflow")
	}
	return nextID
}
