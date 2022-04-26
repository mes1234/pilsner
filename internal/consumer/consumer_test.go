package consumer_test

import (
	consumer "pilsner/internal/consumer"
	"pilsner/internal/stream"
	"testing"
	"time"
)

type streamManagerMock struct {
}

func (s streamManagerMock) Create(streamName string) (err error) {
	return nil
}

func (s streamManagerMock) Delete(streamName string) (err error) {
	return nil
}

func (s streamManagerMock) GetConsumerDataSource(streamName string) (err error, streamIterator <-chan stream.Item) {
	return nil, make(chan stream.Item)
}

func TestBuildingNewConsumer(t *testing.T) {

	manager := consumer.NewMemoryManager(streamManagerMock{})

	err, _ := manager.Create("dummy")

	if err != nil {
		t.Errorf("consumer should be created sucesfully")
	}

}

func TestConsumerReportsConsumedItems(t *testing.T) {

	channel := make(chan stream.Item)

	cons := consumer.NewConsumer(channel)

	count := 4

	for i := 1; i <= count; i++ {
		channel <- stream.Item{}
	}

	time.Sleep(1 * time.Second)

	if cons.ConsumedItems != count {
		t.Errorf("Consummer should consume one item")
	}

}
