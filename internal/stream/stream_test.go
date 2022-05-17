package stream_test

import (
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/stream"
	"testing"
	"time"
)

func TestConsumerOnline(t *testing.T) {

	myStream, cancelStream := stream.NewStream()

	delegate1 := communication.NewDelegate("first")

	counter := 0

	go publishDataToStream(10000, myStream, time.Microsecond, &counter)

	// attach new delegate to stream
	myStream.Stream(*delegate1)

	time.Sleep(5 * time.Second)

	Verify(t, *delegate1)

	delegate1.Cancel()
	cancelStream()

	time.Sleep(1 * time.Second)

}

func TestConsumerLateAttach(t *testing.T) {

	myStream, cancelStream := stream.NewStream()

	delegate1 := communication.NewDelegate("first")

	counter := 0

	go publishDataToStream(10000, myStream, time.Microsecond, &counter)

	time.Sleep(2 * time.Second)
	// attach new delegate to stream
	myStream.Stream(*delegate1)

	time.Sleep(5 * time.Second)

	Verify(t, *delegate1)

	delegate1.Cancel()
	cancelStream()

	time.Sleep(1 * time.Second)
}

func TestConsumerRemoved(t *testing.T) {

	myStream, cancelStream := stream.NewStream()

	delegate1 := communication.NewDelegate("first")

	counter := 0

	go publishDataToStream(10000, myStream, time.Microsecond, &counter)

	time.Sleep(2 * time.Second)
	// attach new delegate to stream
	myStream.Stream(*delegate1)

	time.Sleep(1 * time.Second)

	delegate1.Cancel()

	time.Sleep(1 * time.Second)

	log.Printf("Processing %d", counter)

	Verify(t, *delegate1)

	cancelStream()

	time.Sleep(1 * time.Second)
}

func Verify(t *testing.T, sub1 communication.Delegate) {
	var i = 0
	var finishedFlag = false

	for {
		time.Sleep(time.Millisecond)
		if finishedFlag {
			break
		}
		select {
		case item := <-sub1.Channel:
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
		err := newStream.Publish(communication.Item{
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
