package stream

import (
	"log"
	"sync"
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
	// TODO make implementation which will be non blocking always
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
	subscribers []Delegate
}

func (s *stream) Start(delegate Delegate) {

	var wg sync.WaitGroup

	wg.Add(1)

	go s.read(delegate, &wg)

	s.attach(delegate, &wg)

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
	wg.Wait()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.subscribers = append(s.subscribers, delegate)
	log.Printf("Successfully attached delegate: %s to new data in stream \n", delegate.name)
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
		subs.channel <- item
		log.Printf("Successfully broadcasted to delegate: %s new data %d \n", subs.name, item.Id)
	}
}
