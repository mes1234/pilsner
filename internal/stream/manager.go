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
	GetIterator(streamName string) (err error, streamIterator <-chan Item)
}

func (m *memoryManager) Create(name string) (err error) {
	if _, ok := m.streams.Streams[name]; !ok {
		m.streams.Streams[name] = newStreamer(context{})
		err = nil
		return
	} else {
		err = fmt.Errorf("stream %s already exists", name)
	}
	err = nil
	return
}

func (m *memoryManager) Delete(name string) (err error) {
	return fmt.Errorf("stream %s cannot be deleted", name)
}

func (m *memoryManager) GetIterator(name string) (err error, streamIterator <-chan Item) {
	if streamer, ok := m.streams.Streams[name]; ok {

		// Create channel between source and sink
		channel := make(chan Item, 0)

		// Start writing to channel
		go streamer.Start(channel)

		// return other part of channel to consumer
		streamIterator = channel

		return
	} else {
		err = fmt.Errorf("stream %s doesn't exist", name)
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
