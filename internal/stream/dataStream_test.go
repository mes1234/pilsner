package stream_test

import (
	"context"
	"pilsner/internal/communication"
	"pilsner/internal/stream"
	"testing"
)

func TestDataSourcePutGet(t *testing.T) {

	itemSource := "test"

	ds := stream.NewDataSource(context.Background())

	_ = ds.GetIterator(context.Background())

	ds.Put(communication.Item{
		Source:  itemSource,
		Content: struct{}{},
	})

	err, item := ds.TryGet(0)

	if err != nil {
		t.Errorf("there should be no error in normal operation")
	}

	if item.Source != itemSource {
		t.Errorf("item should have source: %s but got : %s", itemSource, item.Source)
	}
}
func TestDataSourceGetNoPut(t *testing.T) {

	ds := stream.NewDataSource(context.Background())

	_ = ds.GetIterator(context.Background())

	err, item := ds.TryGet(0)

	if err != nil {
		t.Errorf("there should be error ")
	}

	if item.Id != stream.NoItemId {
		t.Errorf("Expected no item ID recieved %d", item.Id)
	}

}
