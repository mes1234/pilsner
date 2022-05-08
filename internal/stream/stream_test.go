package stream_test

import (
	"log"
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestCreatingAndConsumingStream(t *testing.T) {

	myStream := stream.NewStream()

	sub1 := make(chan stream.Item, 100)
	sub2 := make(chan stream.Item, 100)

	delegate1 := stream.NewDelegate(sub1, "first")
	delegate2 := stream.NewDelegate(sub2, "second")

	counter := 0

	go publishDataToStream(10, myStream, &counter)

	// attach new delegate to stream
	myStream.Start(*delegate1)

	time.Sleep(2 * time.Second)

	// attach new delegate to stream
	myStream.Start(*delegate2)

	time.Sleep(2 * time.Second)

	for i := 0; i < counter; i++ {
		item := <-sub1
		if item.Id != i {
			t.Errorf("for sub1 got item with counter %d but expected %d", item.Id, i)
		}
	}

	for i := 0; i < counter; i++ {
		item := <-sub2
		if item.Id != i {
			t.Errorf("for sub2 got item with counter %d but expected %d", item.Id, i)
		}
	}

	log.Printf("Processed %d", counter)
	time.Sleep(2 * time.Second)

	sub3 := make(chan stream.Item, 100)
	delegate3 := stream.NewDelegate(sub3, "third")
	// attach new delegate to stream
	myStream.Start(*delegate3)

	time.Sleep(2 * time.Nanosecond)

	for i := 0; i < counter; i++ {
		item := <-sub3
		if item.Id != i {
			t.Errorf("for sub3 got item with counter %d but expected %d", item.Id, i)
		}
	}

	myStream.Terminate()
}

func publishDataToStream(count int, newStream stream.Publisher, counter *int) {
	for i := 0; i < count; i++ {
		newStream.Publish(stream.Item{
			Id: *counter,
		})
		log.Printf("Added to stream %d item", *counter)
		time.Sleep(time.Nanosecond)
		*counter++
	}
}
