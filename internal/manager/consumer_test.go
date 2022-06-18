package manager_test

import (
	"pilsner/internal/communication"
	"pilsner/internal/manager"
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestAttachedConsumerGetsData(t *testing.T) {

	c := manager.NewConsumerManager()

	st := stream.Get()

	st.Publish(communication.Item{
		Content: "hello",
	})

	_, d := c.Attach("dummy")

	time.Sleep(1)

	data := <-d.Channel

	if data.Id != 0 {
		t.Errorf("Expected 0 got :%d", data.Id)
	}

	if data.Content != "hello" {
		t.Errorf("Expected hello content got :%s", data.Content)
	}

}
