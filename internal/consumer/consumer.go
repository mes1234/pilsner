package consumer

import (
	"fmt"
	"pilsner/internal/stream"
)

type consumer struct {
	streamIterator <-chan stream.Item
	ConsumedItems  int
}

func (c *consumer) startProcessing() {
	for range c.streamIterator {
		c.ConsumedItems++
		fmt.Println("got one")
	}
}

func NewConsumer(stream <-chan stream.Item) *consumer {
	consumer := consumer{
		streamIterator: stream,
	}

	go consumer.startProcessing()

	return &consumer
}
