package consumer

import (
	"github.com/google/uuid"
	. "pilsner/internal/filter"
	. "pilsner/internal/stream"
)

type Consumers struct {
	Consumers map[uuid.UUID]*consumer
}

type consumer struct {
	Stream  <-chan Item
	Filters []*Filter
}

func NewConsumer(stream <-chan Item) *consumer {
	return &consumer{
		Stream:  stream,
		Filters: make([]*Filter, 0),
	}
}
