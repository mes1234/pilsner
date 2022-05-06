package stream_test

import (
	"pilsner/internal/stream"
	"testing"
)

func TestCreatingAndConsumingStream(t *testing.T) {

	newStream := stream.NewStream()

	sub := make(chan stream.Item, 100)

	delegate := stream.NewDelegate(sub, "first")

	for i := 1; i < 5; i++ {
		newStream.Publish(stream.Item{
			Id: i,
		})
	}

	newStream.Start(*delegate)

	for i := 5; i < 10; i++ {
		newStream.Publish(stream.Item{Id: i})
	}

	//for i := 1; i < 10; i++ {
	//	result := <-sub
	//	fmt.Printf("got 1 item %d\n", result.Id)
	//}

	delegate2 := stream.NewDelegate(sub, "second")

	newStream.Start(*delegate2)

	for i := 10; i < 20; i++ {
		newStream.Publish(stream.Item{Id: i})
	}

	//if result == nil {
	//	t.Errorf("got %q, wanted %q", got, want)
	//}
}
