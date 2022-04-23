package stream_test

import (
	"pilsner/internal/stream"
	"testing"
)

func TestBuildingNewStreamManager(t *testing.T) {
	manager := stream.NewMemoryManager()

	_ = manager.Create("dummy", stream.ProtoDefinition{})

	err, strm := manager.Get("dummy")

	if err != nil {
		t.Errorf("There should be no errors")
	}

	if &strm.Proto == nil {
		t.Errorf("There should be protobuf defined")
	}

	if &strm.Items == nil {
		t.Errorf("There should be bootstraped items collection")
	}

}
