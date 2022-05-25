package manager_test

import (
	"pilsner/internal/manager"
	"testing"
)

func TestOnlyOneInstanceOfManagerExists(t *testing.T) {
	streamManager1 := manager.NewStreamManager()
	streamManager2 := manager.NewStreamManager()

	if streamManager1 != streamManager2 {
		t.Errorf("Expected to get the same manager recieved different")
	}

}
