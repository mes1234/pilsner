package stream

import "sync"

const BufferSize int = 100

type Streamer interface {
	RegisterConsumer(channel chan<- Item)
	RegisterPublisher(channel <-chan Item)
}

type stream struct {
	items         []Item
	hash          hash
	context       Context
	buffer        chan Item
	lock          sync.Mutex
	notifications []chan<- Item
}

func (s *stream) RegisterConsumer(channel chan<- Item) {
	s.notifications = append(s.notifications, channel)
}

func (s *stream) RegisterPublisher(channel <-chan Item) {
	for item := range channel {
		s.buffer <- item
	}
}

// startStream begins processing of stream
func (s *stream) startStream() {
	go s.startCollecting()
}

func (s *stream) startCollecting() {
	for item := range s.buffer {
		s.add(item)
		go s.broadcast(item)
	}
}

func (s *stream) broadcast(item Item) {
	for _, channel := range s.notifications {
		channel <- item
	}
}

func (s *stream) add(item Item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, item)
}

func NewStream(context Context) *stream {
	items := make([]Item, 0)
	newStream := stream{
		items:   items,
		context: context,
		buffer:  make(chan Item, BufferSize),
	}
	newStream.startStream()

	return &newStream
}

// Item is a single portion of data in stream
type Item struct {
	content []byte //Raw content of item
}

// Context is container for all metadata regarding stream
type Context struct {
}
