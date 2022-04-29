package stream_test

import (
	"github.com/google/uuid"
	"pilsner/internal/stream"
	"testing"
)

func TestBuildingNewStreamManager(t *testing.T) {
	manager := stream.NewMemoryManager()

	_ = manager.Create("dummy")

	err, streamIterator := manager.CreateConsumerDataSource("dummy", uuid.UUID{})

	if err != nil {
		t.Errorf("There should be no errors")
	}
	if streamIterator == nil {
		t.Errorf("There should be open iterator")
	}

}