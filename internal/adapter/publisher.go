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
	Handle(communication.Item, string, streamManager.StreamManager) error
}

func NewPublisherHandler() *publisherHandler {
	return &publisherHandler{}
}

func (p *publisherHandler) Handle(item communication.Item, streamName string, stream streamManager.StreamManager) error {

	err, publisher := stream.Get(streamName)

	if err != nil {
		return fmt.Errorf("cannot get stream %s", streamName)
	}

	err = publisher.Publish(item)

	if err != nil {
		return fmt.Errorf("failed to publish to stream %s", streamName)
	}

	log.Printf("Published to stream %s item", streamName)

	return nil
}
