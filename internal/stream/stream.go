package stream

import (
	"context"
	"log"
	"sync"
	"time"
)

type Publisher interface {
	Publish(item Item) error
}

type Streamer interface {
	Stream(Delegate)
}

func NewStream(ctx context.Context) *stream {
	buffer := make(chan Item, BufferSize)
	items := NewDataSource()
	newStream := stream{
		items:      items,
		buffer:     buffer,
		terminator: ctx,
	}

	go newStream.run(newStream.terminator)

	return &newStream

}

func (s *stream) Publish(item Item) error {
	// TODO make implementation which will be nonblocking always
	s.buffer <- item
	return nil
}

type stream struct {
	items      Data
	lock       sync.Mutex
	buffer     chan Item
	terminator context.Context
}

func (s *stream) Stream(delegate Delegate) {
	go s.read(delegate, s.items.GetIterator(s.terminator))
}

func (s *stream) read(delegate Delegate, iterator Iterator) {
	log.Printf("Stream publishing historical data for delegate: %s\n", delegate.name)

	for next, item := iterator.Next(); next; next, item = iterator.Next() {
		delegate.channel <- item
	}

	log.Printf("Completed historical data for delegate: %s\n", delegate.name)
}

func (s *stream) run(terminate context.Context) {
	for {
		select {
		case <-terminate.Done():
			log.Printf("Terminating stream operation")
			return
		case item := <-s.buffer:
			s.items.Put(item)
		case <-time.After(Delay):
		}
	}

}
