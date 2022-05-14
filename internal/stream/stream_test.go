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

	Verify(t, sub1)

	cancelConsumer()
	cancelStream()

	time.Sleep(1 * time.Second)

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

	Verify(t, sub1)

	cancelConsumer()
	cancelStream()

	time.Sleep(1 * time.Second)
}

func TestConsumerRemoved(t *testing.T) {

	ctx, cancelStream := context.WithCancel(context.Background())

	myStream := stream.NewStream(ctx)
	sub1 := make(chan stream.Item, 100)
	delegate1 := stream.NewDelegate(sub1, "first")

	counter := 0

	go publishDataToStream(10000, myStream, time.Microsecond, &counter)

	time.Sleep(2 * time.Second)
	// attach new delegate to stream
	cancelConsumer := myStream.Stream(*delegate1)

	time.Sleep(1 * time.Second)

	cancelConsumer()

	time.Sleep(1 * time.Second)

	log.Printf("Processing %d", counter)

	Verify(t, sub1)

	cancelStream()

	time.Sleep(1 * time.Second)
}

func Verify(t *testing.T, sub1 chan stream.Item) {
	var i = 0
	var finishedFlag = false

	for {
		time.Sleep(time.Millisecond)
		if finishedFlag {
			break
		}
		select {
		case item := <-sub1:
			//log.Printf("Got %d from %s expected %d", item.Id, item.Source, i)
			if item.Id != i {
				t.Errorf("for sub1 got item with counter %d but expected %d type : %s", item.Id, i, item.Source)
			}
			i++
		case <-time.After(2 * time.Second):
			finishedFlag = true
		}
	}
}

func publishDataToStream(count int, newStream stream.Publisher, delay time.Duration, counter *int) {
	for i := 0; i < count; i++ {
		err := newStream.Publish(stream.Item{
			Id: *counter,
			Content: []byte{
				0x01,
				0x02,
				0x03,
			},
			Source: "Publisher",
		})
		if err != nil {
			return
		}
		time.Sleep(delay)
		*counter++
	}
}
