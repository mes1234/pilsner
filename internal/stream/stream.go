package stream

import (
	"sync"
)

/////////////
// PUBLIC  //
/////////////

const BufferSize int = 100

type Item struct {
	content []byte
}

type Publisher interface {
	Publish(item Item) error
}

type Iterator interface {
	Start(chan<- Item)
}

func NewStream() *stream {
	buffer := make(chan Item, BufferSize)
	items := make([]Item, BufferSize)
	newStream := stream{
		items:  items,
		buffer: buffer,
	}

	go newStream.run()

	return &newStream

}

func (s *stream) Publish(item Item) error {
	// TODO make implementation which will be non blocking always
	s.buffer <- item
	return nil
}

/////////////
// PRIVATE //
/////////////

type stream struct {
	items  []Item
	buffer chan Item
	lock   sync.Mutex
}

func (s *stream) Start(items chan<- Item) {
	//TODO implement me
	panic("implement me")
}

// run starts operation of stream
func (s *stream) run() {
	for item := range s.buffer {
		s.add(item)
	}
}

// add item to items stored in stream
func (s *stream) add(item Item) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = append(s.items, item)
}
