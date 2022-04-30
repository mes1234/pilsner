package stream

import (
	"fmt"
	"github.com/google/uuid"
)

// Streams contain all streams in pilsner
type Streams struct {
	streams map[string]ConsumerRegistrar
}

type Creator interface {
	Create(streamName string) (err error)
}

type Deleter interface {
	Delete(streamName string) (err error)
}

type DataSourceProvider interface {
	CreateConsumerDataSource(streamName string, consumerId uuid.UUID) (err error, streamIterator <-chan Item)
}

func (m *memoryManager) Create(streamName string) (err error) {
	if _, ok := m.streams.streams[streamName]; !ok {
		m.streams.streams[streamName] = NewStream(Context{})
		err = nil
		return
	} else {
		err = fmt.Errorf("memoryStream %s already exists", streamName)
	}
	err = nil
	return
}

func (m *memoryManager) Delete(streamName string) (err error) {
	return fmt.Errorf("memoryStream %s cannot be deleted", streamName)
}

func (m *memoryManager) CreateConsumerDataSource(streamName string, consumerId uuid.UUID) (err error, streamIterator <-chan Item) {
	if streamer, ok := m.streams.streams[streamName]; ok {

		// Create channel between source and sink
		channel := make(chan Item, 0)

		// Register channel for consumer
		m.consumerChannels[consumerId] = channel

		// StartStreaming writing to channel
		streamer.RegisterConsumer(channel)

		// return other part of channel to consumer
		streamIterator = channel

		return
	} else {
		err = fmt.Errorf("memoryStream %s doesn't exist", streamName)
		return
	}
}

type memoryManager struct {
	streams          Streams
	consumerChannels map[uuid.UUID]chan Item
}

func NewMemoryManager() *memoryManager {
	streams := Streams{
		streams: make(map[string]ConsumerRegistrar),
	}
	return &memoryManager{
		streams:          streams,
		consumerChannels: make(map[uuid.UUID]chan Item),
	}
}
