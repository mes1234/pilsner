package adapter_test

import (
	"pilsner/internal/adapter"
	"pilsner/internal/communication"
	"pilsner/internal/stream"
	"testing"
)

type dummyStreamManager struct {
}

type dummyStreamPublisher struct {
}

func (d *dummyStreamPublisher) Publish(item communication.Item) error {
	return nil
}

func (d *dummyStreamPublisher) Stream(delegate communication.Delegate) {
	//TODO implement me
	panic("implement me")
}

func (d *dummyStreamManager) Add(name string) error {
	//TODO implement me
	panic("implement me")
}

func (d *dummyStreamManager) Remove(name string) error {
	//TODO implement me
	panic("implement me")
}

func (d *dummyStreamManager) Get(name string) (error, stream.StreamerPublisher) {
	return nil, &dummyStreamPublisher{}
}

func TestCorrectHandlerShouldSucceed(t *testing.T) {
	h := adapter.NewPublisherHandler()

	err := h.Handle(communication.Item{}, "dummy", &dummyStreamManager{})

	if err != nil {
		t.Errorf("Expected no error during publishing to stream")
	}
}
