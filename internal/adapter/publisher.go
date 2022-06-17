package adapter

import (
	"fmt"
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/manager/streamManager"
)

type publisherHandler struct {
}

type PublisherHandler interface {
	Handle(communication.Item, streamManager.StreamManager) error
}

func NewPublisherHandler() *publisherHandler {
	return &publisherHandler{}
}

func (p *publisherHandler) Handle(item communication.Item, stream streamManager.StreamManager) error {

	err, publisher := stream.Get()

	if err != nil {
		return fmt.Errorf("cannot get stream ")
	}

	err = publisher.Publish(item)

	if err != nil {
		return fmt.Errorf("failed to publish to stream ")
	}

	log.Printf("Published to stream  item")

	return nil
}
