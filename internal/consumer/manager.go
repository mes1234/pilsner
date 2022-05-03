package consumer

import (
	"fmt"
	"github.com/google/uuid"
	"pilsner/internal/stream"
)

type Creator interface {
	Create(callback Callback) (err error, consumerId uuid.UUID)
}

func (m *memoryManager) Create(callback Callback) (err error, consumerId uuid.UUID) {

	consumerId, _ = uuid.NewUUID()

	if er, dataSource := m.dataSourceProvider.CreateConsumerDataSource(consumerId); er == nil {

		m.consumers[consumerId] = NewConsumer(dataSource, callback)

		return
	} else {
		err = fmt.Errorf("error during creating new consumer for %s", m.dataSourceProvider)
		return
	}
}

type memoryManager struct {
	consumers          map[uuid.UUID]interface{}
	dataSourceProvider stream.DataSourceProvider
}

func NewManager(streamManager stream.DataSourceProvider) *memoryManager {
	consumers := make(map[uuid.UUID]interface{})
	return &memoryManager{
		consumers:          consumers,
		dataSourceProvider: streamManager,
	}
}
