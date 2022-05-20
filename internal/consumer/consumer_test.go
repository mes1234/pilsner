package consumer_test

import (
	"pilsner/internal/communication"
	"pilsner/internal/consumer"
	"testing"
	"time"
)

func TestConsumerConsumesItems(t *testing.T) {

	delegate := communication.NewDelegate("dummy")

	flag := false

	consumer.NewConsumer(*delegate, BuildCallback(&flag), struct{}{})

	time.Sleep(2 * time.Second)

	delegate.Channel <- communication.Item{}

	time.Sleep(2 * time.Second)
	delegate.Cancel()
	if !flag {
		t.Errorf("consumer should consume item and set flag to true")
	}
}

func BuildCallback(flag *bool) consumer.Callback {
	return func(item communication.Item) error {
		*flag = true

		return nil

	}
}
