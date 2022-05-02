package stream

import (
	"fmt"
)

type Creator interface {
	Create(streamName string) (err error)
}

type Deleter interface {
	Delete(streamName string) (err error)
}

type Getter interface {
	Get(streamName string) (err error, stream interface{})
}

func (m *memoryManager) Create(streamName string) (err error) {
	if _, ok := m.streams[streamName]; !ok {
		m.streams[streamName] = NewStream(Context{})
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

type memoryManager struct {
	streams map[string]DataSourceProvider
}

func (m *memoryManager) Get(streamName string) (err error, stream interface{}) {
	if stream, ok := m.streams[streamName]; ok {
		err = nil
		return nil, stream
	} else {
		err = fmt.Errorf("memoryStream %s doesn't exists", streamName)
		return err, nil
	}
}

func NewMemoryManager() *memoryManager {
	return &memoryManager{
		streams: make(map[string]DataSourceProvider),
	}
}
