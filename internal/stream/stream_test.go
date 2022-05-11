package stream_test

import (
	"context"
	"log"
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestConsumerOnline(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	myStream := stream.NewStream(ctx)
	sub1 := make(chan stream.Item, 10)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(2000, myStream, time.Nanosecond, &counter)

	// attach new delegate to stream
	myStream.Stream(*delegate1)

	time.Sleep(10 * time.Second)

	log.Printf("Processing %d", counter)

	for i := 0; i < counter; i++ {
		item := <-sub1
		log.Printf("Got %d from %s expected %d", item.Id, item.Source, i)
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d", item.Id, i)
		}
	}
	cancel()
}

func publishDataToStream(count int, newStream stream.Publisher, delay time.Duration, counter *int) {
	for i := 0; i < count; i++ {
		newStream.Publish(stream.Item{
			Id: *counter,
		})
		//time.Sleep(delay)
		*counter++
	}
}
