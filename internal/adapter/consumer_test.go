package adapter_test

import (
	"context"
	"fmt"
	"pilsner/internal/adapter"
	"pilsner/internal/communication"
	"testing"
	"time"
)

type dummyConsumerManager struct {
	Context context.Context
	Channel chan communication.Item
}

func newDummyConsumerManager(context context.Context, channel chan communication.Item) *dummyConsumerManager {
	return &dummyConsumerManager{
		Context: context,
		Channel: channel,
	}
}

func (d dummyConsumerManager) Attach(streamName string, consumerName string) (error, communication.Delegate) {
	return nil, communication.Delegate{
		Channel: d.Channel,
		Context: d.Context,
	}
}

func TestSendToConsumerAndExit(t *testing.T) {
	// Arrange

	// Item of Interest
	h := adapter.NewConsumerHandler()
	ctx, _ := context.WithCancel(context.Background())
	ch := make(chan communication.Item, 1)

	manager := newDummyConsumerManager(ctx, ch)

	flag := false

	sendFunction := func(flag *bool) adapter.SendFunction {
		return func(item *communication.Item) error {
			*flag = true

			return nil
		}
	}(&flag)

	// Act

	h.HandleSetup(communication.ConsumerSetup{}, manager)

	go h.SendToConsumer(sendFunction)

	ch <- communication.Item{
		Content: []byte{01},
	}

	time.Sleep(1 * time.Second)

	// Assert

	if !flag {
		t.Errorf("Expected flag to be true got %t", flag)
	}
}

type Nothing interface {
}

func TestListenToConsumer(t *testing.T) {
	// Arrange
	flag := false
	count := 0

	receiveFunc := func(m interface{}) error {
		return nil
	}

	handleMsgFunc := func(msg interface{}) error {
		if count < 2 {
			count++
			return nil
		} else {
			flag = true
			return fmt.Errorf("finished")
		}
	}

	// Act

	// Item of Interest
	adapter.ListenToConsumer[Nothing](receiveFunc, handleMsgFunc)

	// Assert

	if flag != true {
		t.Errorf("Expected flag to be true got %t", flag)
	}

}

func TestHandleAck(t *testing.T) {
	// Arrange
	h := adapter.NewConsumerHandler()
	waiter := make(chan struct{}, 1)
	h.Waiter = waiter
	// Act
	h.HandleAck(communication.ConsumerAck{})
	// Assert
	data := <-waiter
	if data != struct{}{} {
		t.Errorf("Expected acknowaldge from waiter got nothin")
	}
}
