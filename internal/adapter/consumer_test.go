package adapter_test

import (
	"context"
	"pilsner/internal/adapter"
	"pilsner/internal/communication"
	"testing"
	"time"
)

func TestSendToConsumerAndExit(t *testing.T) {
	h := adapter.NewConsumerHandler()

	ctx, _ := context.WithCancel(context.Background())

	h.Ctx = ctx

	ch := make(chan communication.Item, 1)

	h.Channel = ch

	flag := false

	sendFunction := func(flag *bool) adapter.SendFunction {
		return func(item *communication.Item) error {
			*flag = true

			return nil
		}
	}(&flag)

	go h.SendToConsumer(sendFunction)

	ch <- communication.Item{
		Content: []byte{01},
	}

	time.Sleep(1 * time.Second)

	if !flag {
		t.Errorf("Expected flag to be true got %t", flag)
	}

}
