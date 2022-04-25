package consumer

import (
	"fmt"
	"github.com/google/uuid"
	"pilsner/internal/stream"
)

type Manager interface {
	Create(streamName string) (err error)
	Delete(streamName string) (err error)
}

func (m *memoryManager) Create(streamName string) (err error, id uuid.UUID) {

	if er, streamIterator := m.streamManager.GetIterator(streamName); er == nil {
		id, _ = uuid.NewUUID()

		m.consumers[id] = newConsumer(streamIterator)

		return
	} else {
		err = fmt.Errorf("error during creating new consumer for %s", streamName)
		return
	}

}

func (m *memoryManager) Delete(streamName string) (err error) {
	//TODO implement me
	panic("implement me")
}

type memoryManager struct {
	consumers     map[uuid.UUID]*consumer
	streamManager stream.Manager
}

func NewMemoryManager(streamManager stream.Manager) *memoryManager {
	consumers := make(map[uuid.UUID]*consumer)
	return &memoryManager{
		consumers:     consumers,
		streamManager: streamManager,
	}
}
