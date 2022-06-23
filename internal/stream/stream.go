package stream

import (
	"context"
	"fmt"
	"log"
	"pilsner/internal/communication"
	"pilsner/setup"
	"sync"
)

var initStream sync.Once

// memoryStream is app instance of stream
var memoryStream StreamerPublisher

// cancel is used to gracefully shutdown stream and app
var cancel context.CancelFunc

func init() {
	initStream.Do(func() {
		memoryStream, cancel = NewStream()
	})
}

func Get() StreamerPublisher {
	return memoryStream
}

func Shutdown() {
	cancel()
}

type StreamerPublisher interface {
	Publisher
	Streamer
}

type Publisher interface {
	// Publish writes single item to stream
	Publish(item communication.Item) error
}

type Streamer interface {
	// Stream starts streaming given stream to delegate
	//
	// Streaming can be cancelled using returned CancelFunc
	Stream(communication.Delegate)
}

func NewStream() (*stream, context.CancelFunc) {

	ctx, cancelStream := context.WithCancel(context.Background())

	buffer := make(chan communication.Item, setup.Config.BufferSize)
	items := NewDataSource(ctx)
	newStream := stream{
		items:      items,
		buffer:     buffer,
		terminator: ctx,
	}

	return &newStream, cancelStream
}

func (s *stream) Publish(item communication.Item) error {
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
	lock       sync.RWMutex
	buffer     chan communication.Item
	terminator context.Context
}

func (s *stream) Stream(delegate communication.Delegate) {

	go s.read(delegate, s.items.GetIterator(delegate.Context))
}

func (s *stream) read(delegate communication.Delegate, iterator Iterator) {
	log.Printf("Stream publishing historical data for delegate: %s\n", delegate.Name)

	for next, item := iterator.Next(); next; next, item = iterator.Next() {
		delegate.Channel <- item
	}

	log.Printf("Completed historical data for delegate: %s\n", delegate.Name)
}
