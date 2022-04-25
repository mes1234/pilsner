package consumer_test

import (
	consumer "pilsner/internal/consumer"
	"pilsner/internal/stream"
	"testing"
)

type streamManagerMock struct {
}

func (s streamManagerMock) Create(streamName string) (err error) {
	return nil
}

func (s streamManagerMock) Delete(streamName string) (err error) {
	return nil
}

func (s streamManagerMock) GetIterator(streamName string) (err error, streamIterator <-chan stream.Item) {
	return nil, make(chan stream.Item)
}

func TestBuildingNewConsumer(t *testing.T) {

	manager := consumer.NewMemoryManager(streamManagerMock{})

	err, _ := manager.Create("dummy")

	if err != nil {
		t.Errorf("consumer should be created sucesfully")
	}

}
