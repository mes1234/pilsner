package stream

import (
	"log"
	"sync"
	"time"
)

/////////////
// PUBLIC  //
/////////////

const BufferSize int = 100

type Item struct {
	Id      int
	content []byte
}

type Publisher interface {
	Publish(item Item) error
}

type Iterator interface {
	Start(Delegate)
}

type Delegate struct {
	channel chan<- Item
	name    string
}

func NewDelegate(channel chan<- Item, name string) *Delegate {
	return &Delegate{
		name:    name,
		channel: channel,
	}
}

func NewStream() *stream {
	buffer := make(chan Item, BufferSize)
	items := make([]Item, 0)
	newStream := stream{
		items:  items,
		buffer: buffer,
	}

	go newStream.run()

	return &newStream

}

func (s *stream) Publish(item Item) error {
	// TODO make implementation which will be nonblocking always
	s.buffer <- item
	return nil
}

/////////////
// PRIVATE //
/////////////

type stream struct {
	items       []Item
	buffer      chan Item
	lock        sync.Mutex
	position    int
	subscribers []chan<- Item
}

func (s *stream) Start(delegate Delegate) {

	// wg will indicate when sync up end and new data can be streamed
	wg := buildWaitGroup()

	// attach delegate to new data
	s.attach(delegate, wg)

	// begin sync up
	go s.read(delegate, wg)

}

func buildWaitGroup() *sync.WaitGroup {
	var wg sync.WaitGroup

	wg.Add(1)
	return &wg
}

func (s *stream) read(delegate Delegate, wg *sync.WaitGroup) {
	log.Printf("Start publishing historical data for delegate: %s\n", delegate.name)
	for _, item := range s.items {
		delegate.channel <- item
		log.Printf("Published historical data for delegate %s item: %d\n", delegate.name, item.Id)
	}
	wg.Done()
	log.Printf("Completed publishing historical data for delegate: %s\n", delegate.name)
}

// run starts operation of stream
func (s *stream) run() {
	for item := range s.buffer {
		s.add(item)
	}
}

func (s *stream) attach(delegate Delegate, wg *sync.WaitGroup) {

	log.Printf("Requested to attach delegate: %s to new data in stream \n", delegate.name)
	s.lock.Lock()
	defer s.lock.Unlock()

	// Create sync time buffer
	buf := make(chan Item, BufferSize)
	// Attach buffer to broadcaster so new data is queued
	s.subscribers = append(s.subscribers, buf)
	log.Printf("Successfully attached delegate: %s to new data in stream \n", delegate.name)

	// start shuffling data between buffer and delegate
	go shuffle(buf, delegate.channel, wg)
}

// shuffle data from source to destination when waiter is ready
func shuffle(source <-chan Item, destination chan<- Item, waiter *sync.WaitGroup) {
	waiter.Wait()

	for {
		select {
		case msg := <-source:
			select {
			case destination <- msg:
				log.Printf("Successfully shuffled data to: <> id:%d\n", msg.Id)
			case <-time.After(time.Microsecond):
				log.Printf("Buffer full need to wait")
			}

		}
	}

}

// add item to items stored in stream
func (s *stream) add(item Item) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = append(s.items, item)
	s.position = len(s.items)
	s.broadcast(item)
}

// broadcast information about new item in stream to all subscribers
func (s *stream) broadcast(item Item) {
	for _, subs := range s.subscribers {
		subs <- item
		log.Printf("Successfully broadcasted to delegate: %s new data %d \n", "dummy", item.Id)
	}
}
