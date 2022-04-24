package consumer_test

import (
	consumer "pilsner/internal/consumer"
	"pilsner/internal/stream"
	"testing"
)

func TestBuildingNewConsumer(t *testing.T) {

	streamMock := make(<-chan stream.Item)

	con := consumer.NewConsumer(streamMock)

	if con.Filters == nil {
		t.Errorf("consumer should have initialized Filters")
	}

	if con.Stream == nil {
		t.Errorf("consumer shall always be attached to stream")
	}
}
