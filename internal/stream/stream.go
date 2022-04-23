package stream

// Streams contain all streams in pilsner
type Streams struct {
	Streams map[string]Stream //Streams in given instance
}

type Stream struct {
	Items []*Item         // All items in stream
	Proto ProtoDefinition // Assigned protobuf for stream
}

func NewStream(definition ProtoDefinition) *Stream {
	items := make([]*Item, 0)
	return &Stream{
		Items: items,
		Proto: definition,
	}
}

type Item struct {
	content []byte //Raw content of item
}

type ProtoDefinition struct {
}

type Context struct {
}
