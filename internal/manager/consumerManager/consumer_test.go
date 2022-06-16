package consumerManager_test

import (
	"pilsner/internal/communication"
	"pilsner/internal/manager/consumerManager"
	"pilsner/internal/manager/streamManager"
	"testing"
	"time"
)

func TestAttachedConsumerToNonValidStreamIsInvalid(t *testing.T) {

	c := consumerManager.NewConsumerManager()

	err, _ := c.Attach("dummy", "dummy")

	if err == nil {
		t.Errorf("Expected to fail attaching consumer to stream")
	}
}

func TestAttachedConsumerGetsData(t *testing.T) {

	s := streamManager.NewStreamManager()
	c := consumerManager.NewConsumerManager()

	s.Add("dummy")

	_, st := s.Get("dummy")

	st.Publish(communication.Item{
		Content: "hello",
	})

	_, d := c.Attach("dummy", "dummy")

	time.Sleep(1)

	data := <-d.Channel

	if data.Id != 0 {
		t.Errorf("Expected 0 got :%d", data.Id)
	}

	if data.Content != "hello" {
		t.Errorf("Expected hello content got :%s", data.Content)
	}

}
