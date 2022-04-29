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

	id, _ = uuid.NewUUID()

	if er, dataSource := m.streamManager.CreateConsumerDataSource(streamName, id); er == nil {

		m.consumers[id] = NewConsumer(dataSource)

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
	consumers     map[uuid.UUID]Consumer
	streamManager stream.Manager
}

func NewMemoryManager(streamManager stream.Manager) *memoryManager {
	consumers := make(map[uuid.UUID]Consumer)
	return &memoryManager{
		consumers:     consumers,
		streamManager: streamManager,
	}
}
