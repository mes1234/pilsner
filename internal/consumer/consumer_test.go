package consumer_test

import (
	consumer "pilsner/internal/consumer"
	"pilsner/internal/stream"
	"testing"
)

func TestBuildingNewConsumer(t *testing.T) {
	newStream := stream.Stream{}
	con := consumer.NewConsumer(&newStream)

	if con.Filters == nil {
		t.Errorf("consumer should have initialized Filters")
	}

	if con.Stream == nil {
		t.Errorf("consumer shall always be attached to stream")
	}
}
