package stream_test

import (
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestStreamPublishing(t *testing.T) {
	strm := stream.NewStream(stream.Context{})

	pub := make(chan stream.Item)
	sub := make(chan stream.Item, 1)

	go strm.RegisterPublisher(pub)

	strm.RegisterConsumer(sub)

	pub <- stream.Item{}

	time.Sleep(1 * time.Second)

	if len(sub) != cap(sub) {
		t.Errorf("there is no pending data in consumer ")
	}

}
