package stream

import (
	"log"
	"sync"
	"time"
)

const BufferSize int = 1000000

type Publisher interface {
	Publish(item Item) error
}

type Iterator interface {
	Start(Delegate)
}

type Terminator interface {
	Terminate()
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

func (s *stream) Terminate() {
	s.terminator <- struct{}{}
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

func (s *stream) Start(delegate Delegate) {

	// wg will indicate when sync up end and new data can be streamed
	wg := buildWaitGroup()

	// attach delegate to new data
	s.attach(delegate, wg)

	// begin sync up
	go s.read(delegate, wg, s.items.GetPosition())

}

func buildWaitGroup() *sync.WaitGroup {
	var wg sync.WaitGroup

	wg.Add(1)
	return &wg
}

func (s *stream) read(delegate Delegate, wg *sync.WaitGroup, length int) {
	log.Printf("Start publishing historical data for delegate: %s\n", delegate.name)

	for i := 0; i <= length; i++ {
		item := s.items.Get(i)
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
		default:
			time.Sleep(100 * time.Millisecond)
			log.Printf("No data waiting...")
		}
	}

}
