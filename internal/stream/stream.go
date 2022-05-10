package stream

import (
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

func NewStream() *stream {
	buffer := make(chan Item, BufferSize)
	items := NewDataSource()
	newStream := stream{
		items:      items,
		buffer:     buffer,
		terminator: make(chan struct{}),
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
	items       Data
	lock        sync.Mutex
	buffer      chan Item
	subscribers []subscriber
	terminator  chan struct{}
}

func (s *stream) Stream(delegate Delegate) {

	// wg will indicate when sync up end and new data can be streamed
	wg := buildWaitGroup()

	// attach delegate to new data
	s.attach(delegate, wg)

	// begin sync up
	go s.read(delegate, wg, s.items.GetIterator())

}

func buildWaitGroup() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	return &wg
}

func (s *stream) read(delegate Delegate, wg *sync.WaitGroup, iterator Iterator) {
	log.Printf("Stream publishing historical data for delegate: %s\n", delegate.name)

	for next, item := iterator.Next(); next; next, item = iterator.Next() {
		delegate.channel <- item
		log.Printf("Published historical data for delegate %s item: %d\n", delegate.name, item.Id)
	}

	wg.Done()
	log.Printf("Completed publishing historical data for delegate: %s\n", delegate.name)
}

func (s *stream) attach(delegate Delegate, wg *sync.WaitGroup) {

	log.Printf("Requested to attach delegate: %s to new data in stream \n", delegate.name)

	newSubscriber := subscriber{
		channel: delegate,
		wg:      wg,
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	// Attach buffer to broadcaster so new data is queued
	s.subscribers = append(s.subscribers, newSubscriber)
	log.Printf("Successfully attached delegate: %s to new data in stream \n", delegate.name)

}

func (s *stream) run(terminate <-chan struct{}) {
	for {
		select {
		case <-terminate:
			log.Printf("Terminating stream operation")
			return
		case item := <-s.buffer:
			// save item to stream
			s.items.Put(item)
			for _, sub := range s.subscribers {
				sub.shuffle(item)
			}
		case <-time.After(Delay):
			log.Printf("No data waiting...")
		}
	}

}
