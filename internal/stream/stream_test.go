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
	sub1 := make(chan stream.Item, 100)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(10000, myStream, time.Microsecond, &counter)

	// attach new delegate to stream
	cancelConsumer := myStream.Stream(*delegate1)

	time.Sleep(5 * time.Second)

	log.Printf("Processing %d", counter)

	var i = 0

	for ; i < counter; i++ {
		time.Sleep(time.Nanosecond)
		item := <-sub1
		//	log.Printf("Got %d from %s expected %d", item.Id, item.Source, i)
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d type : %s", item.Id, i, item.Source)
		}
	}
	log.Printf("Processed %d", i)
	cancelConsumer()
	cancelStream()

}

func TestConsumerLateAttach(t *testing.T) {

	ctx, cancelStream := context.WithCancel(context.Background())

	myStream := stream.NewStream(ctx)
	sub1 := make(chan stream.Item, 100)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(10000, myStream, time.Microsecond, &counter)

	time.Sleep(2 * time.Second)
	// attach new delegate to stream
	cancelConsumer := myStream.Stream(*delegate1)

	time.Sleep(5 * time.Second)

	log.Printf("Processing %d", counter)

	var i = 0

	for ; i < counter; i++ {
		time.Sleep(time.Nanosecond)
		item := <-sub1
		//log.Printf("Got %d from %s expected %d", item.Id, item.Source, i)
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d type : %s", item.Id, i, item.Source)
		}
	}

	log.Printf("Processed %d", i)

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
		time.Sleep(delay)
		*counter++
	}
}
