package consumer_test

import (
	"fmt"
	"github.com/google/uuid"
	consumer "pilsner/internal/consumer"
	"pilsner/internal/stream"
	"testing"
	"time"
)

type streamManagerMock struct {
}

func (s streamManagerMock) CreateConsumerDataSource(consumerId uuid.UUID) (err error, streamIterator <-chan stream.Item) {
	return nil, make(chan stream.Item)
}

func TestBuildingNewConsumer(t *testing.T) {

	manager := consumer.NewMemoryManager(streamManagerMock{})

	err, _ := manager.Create()

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

func TestConsumerCallback(t *testing.T) {
	channel := make(chan stream.Item)

	cons := consumer.NewConsumer(channel)

	count := 4

	hit := 0

	cons.RegisterCallback(func(item stream.Item) error {
		hit++
		return nil
	})

	for i := 1; i <= count; i++ {
		channel <- stream.Item{}
	}

	time.Sleep(1 * time.Second)

	if hit != count {
		t.Errorf("Consummer should consume items item")
	}
}

func TestConsumerCallbackRetry(t *testing.T) {
	channel := make(chan stream.Item)

	cons := consumer.NewConsumer(channel)

	count := 4

	hit := 0

	cons.RegisterCallback(func(item stream.Item) error {
		hit++
		return fmt.Errorf("random callback error")
	})

	for i := 1; i <= count; i++ {
		channel <- stream.Item{}
	}

	time.Sleep(1 * time.Second)

	if hit != count*consumer.DefaultRetryAttempts {
		t.Errorf("Consummer should consume items item")
	}
}
