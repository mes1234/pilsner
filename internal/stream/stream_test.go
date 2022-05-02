package stream_test

import (
	"github.com/google/uuid"
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestStreamPublishing(t *testing.T) {
	strm := stream.NewStream(stream.Context{})

	pub := make(chan stream.Item)

	strm.RegisterPublisher(pub)

	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()

	_, con1 := strm.CreateConsumerDataSource(id1)
	_, con2 := strm.CreateConsumerDataSource(id2)

	pub <- stream.Item{}

	time.Sleep(1 * time.Second)

	if len(con1) != cap(con1) {
		t.Errorf("there is no pending data in consumer ")
	}

	if len(con2) != cap(con2) {
		t.Errorf("there is no pending data in consumer ")
	}

}
