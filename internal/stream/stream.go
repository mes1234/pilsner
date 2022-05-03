package stream

import (
	"github.com/google/uuid"
	"sync"
)

const BufferSize int = 100

type PublisherRegistrar interface {
	RegisterPublisher(channel <-chan Item)
}

type DataSourceProvider interface {
	CreateConsumerDataSource(consumerId uuid.UUID) (err error, streamIterator <-chan Item)
}

type memoryStream struct {
	items            []Item
	context          Context
	buffer           chan Item
	lock             sync.Mutex
	consumerChannels map[uuid.UUID]chan<- Item
}

func (ms *memoryStream) CreateConsumerDataSource(consumerId uuid.UUID) (err error, streamIterator <-chan Item) {

	// Create channel between source and sink
	channel := make(chan Item, 0)

	// Register channel for consumer
	ms.consumerChannels[consumerId] = channel

	// return other part of channel to consumer
	streamIterator = channel

	return
}

func (ms *memoryStream) RegisterPublisher(channel <-chan Item) {
	go ms.startPublisher(channel)
}

func (ms *memoryStream) startPublisher(channel <-chan Item) {
	for item := range channel {
		ms.buffer <- item
	}
}

func (ms *memoryStream) startPublishing() {
	for item := range ms.buffer {
		ms.add(item)
		go ms.broadcast(item)
	}
}

func (ms *memoryStream) broadcast(item Item) {
	for _, channel := range ms.consumerChannels {
		channel <- item
	}
}

func (ms *memoryStream) add(item Item) {

	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.items = append(ms.items, item)
}

func NewStream(context Context) *memoryStream {
	items := make([]Item, 0)

	newStream := memoryStream{
		items:            items,
		context:          context,
		buffer:           make(chan Item, BufferSize),
		consumerChannels: make(map[uuid.UUID]chan<- Item, BufferSize),
	}

	go newStream.startPublishing()

	return &newStream
}
