package streamManager_test

import (
	"pilsner/internal/manager/streamManager"
	"testing"
)

func TestOnlyOneInstanceOfManagerExists(t *testing.T) {
	streamManager1 := streamManager.NewStreamManager()
	streamManager2 := streamManager.NewStreamManager()

	if streamManager1 != streamManager2 {
		t.Errorf("Expected to get the same manager recieved different")
	}

}

func TestAddedStreamIsAvailableInAllInstances(t *testing.T) {
	streamManager1 := streamManager.NewStreamManager()
	streamManager2 := streamManager.NewStreamManager()

	streamManager1.Add("dummy")

	err, _ := streamManager2.Get("dummy")

	if err != nil {
		t.Errorf("Expected no error while retrieving defined stream")
	}

}

func TestRemovedStreamIsNotAvailableInAllInstances(t *testing.T) {
	streamManager1 := streamManager.NewStreamManager()
	streamManager2 := streamManager.NewStreamManager()

	streamManager1.Add("dummy")

	err, _ := streamManager2.Get("dummy")

	if err != nil {
		t.Errorf("Expected no error while retrieving defined stream")
	}

	streamManager1.Remove("dummy")

	err, _ = streamManager2.Get("dummy")

	if err == nil {
		t.Errorf("Expected  error while retrieving removed stream")
	}

}
