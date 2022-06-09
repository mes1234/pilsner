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

	ctx, cancel := context.WithCancel(context.Background())

	h.Ctx = ctx

	ch := make(chan communication.Item, 1)

	h.Channel = ch

	send := func(flag *bool) adapter.SendFunction {
		return func(item *communication.Item) error {
			*flag = true

			return nil
		}
	}

	flag := false

	go h.SendToConsumer(send(&flag))

	ch <- communication.Item{
		Content: true,
	}

	time.Sleep(100)

	cancel()

	if !flag {
		t.Errorf("Expected flag to be true got %t", flag)
	}

}
