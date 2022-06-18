package adapter_test

import (
	"pilsner/internal/adapter"
	"pilsner/internal/communication"
	"testing"
)

type dummyStream struct {
}

func (d *dummyStream) Publish(item communication.Item) error {
	return nil
}

func TestCorrectHandlerShouldSucceed(t *testing.T) {
	h := adapter.NewPublisherHandler()

	err := h.Handle(communication.Item{}, &dummyStream{})

	if err != nil {
		t.Errorf("Expected no error during publishing to stream")
	}
}
