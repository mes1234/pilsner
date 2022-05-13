package stream

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Publisher interface {
	Publish(item Item) error
}

type Streamer interface {
	Stream(Delegate)
}

func NewStream(ctx context.Context) *stream {
	buffer := make(chan Item, BufferSize)
	items := NewDataSource(ctx)
	newStream := stream{
		items:      items,
		buffer:     buffer,
		terminator: ctx,
	}

	return &newStream
}

func (s *stream) Publish(item Item) error {
	select {
	case <-s.terminator.Done():
		return fmt.Errorf("stream operation ended")
	default:
		s.items.Put(item)
		return nil
	}
}

type stream struct {
	items      Data
	lock       sync.Mutex
	buffer     chan Item
	terminator context.Context
}

func (s *stream) Stream(delegate Delegate) context.CancelFunc {

	ctx, cancel := context.WithCancel(context.Background())

	go s.read(delegate, s.items.GetIterator(ctx))

	return cancel
}

func (s *stream) read(delegate Delegate, iterator Iterator) {
	log.Printf("Stream publishing historical data for delegate: %s\n", delegate.name)

	for next, item := iterator.Next(); next; next, item = iterator.Next() {
		delegate.channel <- item
	}

	log.Printf("Completed historical data for delegate: %s\n", delegate.name)
}
