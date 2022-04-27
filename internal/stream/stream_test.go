package stream_test

import (
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestStreamPublishing(t *testing.T) {
	strm := stream.NewStream(stream.Context{})

	pub := make(chan stream.Item)
	con1 := make(chan stream.Item, 1)
	con2 := make(chan stream.Item, 1)

	strm.RegisterPublisher(pub)

	strm.RegisterConsumer(con1)
	strm.RegisterConsumer(con2)

	pub <- stream.Item{}

	time.Sleep(1 * time.Second)

	if len(con1) != cap(con1) {
		t.Errorf("there is no pending data in consumer ")
	}

	if len(con2) != cap(con2) {
		t.Errorf("there is no pending data in consumer ")
	}

}
