package stream

import "sync"

const BufferSize int = 100

type PublisherRegistrar interface {
	RegisterPublisher(channel <-chan Item)
}

type ConsumerRegistrar interface {
	RegisterConsumer(channel chan<- Item)
}

type memoryStream struct {
	items         []Item
	context       Context
	buffer        chan Item
	lock          sync.Mutex
	notifications []chan<- Item
}

func (s *memoryStream) RegisterConsumer(channel chan<- Item) {
	s.notifications = append(s.notifications, channel)
}

func (s *memoryStream) RegisterPublisher(channel <-chan Item) {
	go s.startPublisher(channel)
}

func (s *memoryStream) startPublisher(channel <-chan Item) {
	for item := range channel {
		s.buffer <- item
	}
}

func (s *memoryStream) startCollecting() {
	for item := range s.buffer {
		s.add(item)
		go s.broadcast(item)
	}
}

func (s *memoryStream) broadcast(item Item) {
	for _, channel := range s.notifications {
		channel <- item
	}
}

func (s *memoryStream) add(item Item) {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = append(s.items, item)
}

func NewStream(context Context) *memoryStream {
	items := make([]Item, 0)

	newStream := memoryStream{
		items:   items,
		context: context,
		buffer:  make(chan Item, BufferSize),
	}

	go newStream.startCollecting()

	return &newStream
}
