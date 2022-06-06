package handler

import (
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/manager/streamManager"
)

type PublisherHandler struct {
}

func NewPublisherHandler() *PublisherHandler {
	return &PublisherHandler{}
}

func (p *PublisherHandler) Handle(item communication.Item, streamName string) error {

	stream := streamManager.NewStreamManager()

	_, publisher := stream.Get(streamName)

	_ = publisher.Publish(item)

	log.Printf("Published to stream %s item", streamName)

	return nil
}
