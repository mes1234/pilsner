package stream_test

import (
	"log"
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestConsumerOnline(t *testing.T) {

	myStream := stream.NewStream()
	sub1 := make(chan stream.Item, 1)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(200000000, myStream, time.Nanosecond, &counter)

	// attach new delegate to stream
	myStream.Stream(*delegate1)

	time.Sleep(10 * time.Second)

	log.Printf("Processing %d", counter)

	for i := 0; i < counter; i++ {
		item := <-sub1
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d", item.Id, i)
		}
	}

	myStream.Terminate()
}

func TestConsumerLateAttach(t *testing.T) {

	myStream := stream.NewStream()
	sub1 := make(chan stream.Item, 1)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(200, myStream, time.Nanosecond, &counter)

	time.Sleep(1 * time.Second)
	// attach new delegate to stream
	myStream.Stream(*delegate1)

	log.Printf("Processing %d", counter)

	for i := 0; i < counter; i++ {
		item := <-sub1
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d", item.Id, i)
		}
	}

	myStream.Terminate()
}

func TestConsumerOnlineAndLaterPublishedMoreData(t *testing.T) {
	myStream := stream.NewStream()
	sub1 := make(chan stream.Item, 10)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(200, myStream, time.Nanosecond, &counter)

	time.Sleep(5 * time.Second)
	// attach new delegate to stream
	myStream.Stream(*delegate1)

	log.Printf("Processing %d", counter)

	i := 0

	for ; i < counter; i++ {
		item := <-sub1
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d", item.Id, i)
		}

	}

	go publishDataToStream(200, myStream, time.Nanosecond, &counter)

	time.Sleep(5 * time.Second)

	log.Printf("Processing %d", counter)

	for ; i < counter; i++ {
		item := <-sub1
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d", item.Id, i)
		}
	}

	myStream.Terminate()
}

func publishDataToStream(count int, newStream stream.Publisher, delay time.Duration, counter *int) {
	for i := 0; i < count; i++ {
		newStream.Publish(stream.Item{
			Id: *counter,
		})
		log.Printf("Added to stream %d item", *counter)
		time.Sleep(delay)
		*counter++
	}
}
