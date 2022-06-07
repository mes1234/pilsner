package adapter

import (
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/manager/streamManager"
)

type publisherHandler struct {
}

type PublisherHandler interface {
	Handle(item communication.Item, streamName string) error
}

func NewPublisherHandler() *publisherHandler {
	return &publisherHandler{}
}

func (p *publisherHandler) Handle(item communication.Item, streamName string) error {

	stream := streamManager.NewStreamManager()

	_, publisher := stream.Get(streamName)

	_ = publisher.Publish(item)

	log.Printf("Published to stream %s item", streamName)

	return nil
}
