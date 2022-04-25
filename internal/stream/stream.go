package stream

type Streamer interface {
	Start(channel chan<- Item)
}

type stream struct {
	items   []*Item
	hash    hash
	context context
}

func (s *stream) Start(channel chan<- Item) {
	// TODO this is not real implementation
	for _, value := range s.items {
		channel <- *value
	}
}

func newStreamer(context context) *stream {
	items := make([]*Item, 0)
	newStream := stream{
		items:   items,
		context: context,
	}
	return &newStream
}

// Item is a single portion of data in stream
type Item struct {
	content []byte //Raw content of item
}

// context is container for all metadata regarding stream
type context struct {
}
