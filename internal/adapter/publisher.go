package adapter

import (
	"fmt"
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/stream"
)

type publisherHandler struct {
}

type PublisherHandler interface {
	Handle(communication.Item, stream.Publisher) error
}

func NewPublisherHandler() *publisherHandler {
	return &publisherHandler{}
}

func (p *publisherHandler) Handle(item communication.Item, stream stream.Publisher) error {

	err := stream.Publish(item)

	if err != nil {
		return fmt.Errorf("failed to publish to stream ")
	}

	log.Printf("Published to stream  item")

	return nil
}
