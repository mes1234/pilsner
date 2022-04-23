package stream

import (
	"fmt"
	"github.com/google/uuid"
)

type Manager interface {
	Create(name string, protobuf ProtoDefinition) (err error, id uuid.UUID)
	Delete(name string) (err error)
	Get(name string) (err error, stream *Stream)
}

func (m *memoryManager) Create(name string, protobuf ProtoDefinition) (err error) {
	if _, ok := m.streams.Streams[name]; !ok {
		m.streams.Streams[name] = *NewStream(protobuf)
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

func (m *memoryManager) Get(name string) (err error, stream *Stream) {
	if val, ok := m.streams.Streams[name]; ok {
		stream = &val
		err = nil
		return
	} else {
		err = fmt.Errorf("stream %s doesn't exist", name)
	}
	err = nil
	return
}

type memoryManager struct {
	streams Streams
}

func NewMemoryManager() *memoryManager {
	streams := Streams{
		Streams: make(map[string]Stream),
	}
	return &memoryManager{streams: streams}
}
