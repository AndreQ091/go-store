package core

const (
	_                  = iota
	EventPut EventType = iota
	EventDelete
)

type TransactionLogger interface {
	WritePut(key, value string)
	WriteDelete(key string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
}

type Event struct {
	Sequence uint64
	Type     EventType
	Key      string
	Value    string
}

type EventType byte