package consumer

import (
	"github.com/google/uuid"
	. "pilsner/internal/filter"
	. "pilsner/internal/stream"
)

type Consumers struct {
	Consumers map[uuid.UUID]*Consumer
}

type Consumer struct {
	Stream  *Stream
	Filters []*Filter
}
