package consumer

import (
	"pilsner/internal/stream"
)

type consumer struct {
	streamIterator <-chan stream.Item
}

func newConsumer(stream <-chan stream.Item) *consumer {
	return &consumer{
		streamIterator: stream,
	}
}
