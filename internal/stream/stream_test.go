package stream_test

import (
	"context"
	"log"
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestConsumerOnline(t *testing.T) {

	ctx, cancelStream := context.WithCancel(context.Background())

	myStream := stream.NewStream(ctx)
	sub1 := make(chan stream.Item, 10)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(2000000, myStream, time.Microsecond, &counter)

	// attach new delegate to stream
	cancelConsumer := myStream.Stream(*delegate1)

	time.Sleep(2 * time.Second)

	log.Printf("Processing %d", counter)

	for i := 0; i < counter; i++ {
		item := <-sub1
		//log.Printf("Got %d from %s expected %d", item.Id, item.Source, i)
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d type : %s", item.Id, i, item.Source)
		}
	}

	cancelConsumer()
	cancelStream()

}

func TestConsumerLateAttach(t *testing.T) {

	ctx, cancelStream := context.WithCancel(context.Background())

	myStream := stream.NewStream(ctx)
	sub1 := make(chan stream.Item, 10)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(3400000, myStream, time.Microsecond, &counter)

	time.Sleep(2 * time.Second)
	// attach new delegate to stream
	cancelConsumer := myStream.Stream(*delegate1)

	time.Sleep(2 * time.Second)

	log.Printf("Processing %d", counter)

	for i := 0; i < counter; i++ {
		item := <-sub1
		//log.Printf("Got %d from %s expected %d", item.Id, item.Source, i)
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d type : %s", item.Id, i, item.Source)
		}
	}

	cancelConsumer()
	cancelStream()

}

func publishDataToStream(count int, newStream stream.Publisher, delay time.Duration, counter *int) {
	for i := 0; i < count; i++ {
		err := newStream.Publish(stream.Item{
			Id: *counter,
		})
		if err != nil {
			return
		}
		//time.Sleep(delay)
		*counter++
	}
}
