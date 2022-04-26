package stream

import (
	"fmt"
)

// Streams contain all streams in pilsner
type Streams struct {
	Streams map[string]Streamer //Streams in given instance
}

type Manager interface {
	Create(streamName string) (err error)
	Delete(streamName string) (err error)
	GetConsumerDataSource(streamName string) (err error, streamIterator <-chan Item)
}

func (m *memoryManager) Create(streamName string) (err error) {
	if _, ok := m.streams.Streams[streamName]; !ok {
		m.streams.Streams[streamName] = NewStream(Context{})
		err = nil
		return
	} else {
		err = fmt.Errorf("stream %s already exists", streamName)
	}
	err = nil
	return
}

func (m *memoryManager) Delete(streamName string) (err error) {
	return fmt.Errorf("stream %s cannot be deleted", streamName)
}

func (m *memoryManager) GetConsumerDataSource(streamName string) (err error, streamIterator <-chan Item) {
	if streamer, ok := m.streams.Streams[streamName]; ok {

		// Create channel between source and sink
		channel := make(chan Item, 0)

		// StartStreaming writing to channel
		streamer.RegisterConsumer(channel)

		// return other part of channel to consumer
		streamIterator = channel

		return
	} else {
		err = fmt.Errorf("stream %s doesn't exist", streamName)
		return
	}
}

type memoryManager struct {
	streams Streams
}

func NewMemoryManager() *memoryManager {
	streams := Streams{
		Streams: make(map[string]Streamer),
	}
	return &memoryManager{streams: streams}
}
