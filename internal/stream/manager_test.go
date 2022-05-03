package stream_test

import (
	"pilsner/internal/stream"
	"testing"
)

func TestBuildingNewStreamManager(t *testing.T) {
	manager := stream.NewManager()

	_ = manager.Create("dummy")

	err, streamIterator := manager.Get("dummy")

	if err != nil {
		t.Errorf("There should be no errors")
	}
	if streamIterator == nil {
		t.Errorf("There should be stream")
	}

}
