package stream

import (
	"github.com/google/uuid"
)

// Streams contain all streams in pilsner
type Streams struct {
	Streams map[uuid.UUID]Stream //
}

type Stream struct {
	Items []*Item         // All items in stream
	Proto ProtoDefinition // Assigned protobuf for stream
}

type Item struct {
	content []byte //Raw content of item
}

type ProtoDefinition struct {
}
