package stream_test

import (
	"context"
	"fmt"
	"math/rand"
	"pilsner/internal/communication"
	"pilsner/internal/stream"
	"testing"
)

type dataMock struct {
	counter        int
	switchPosition int
}

func (d *dataMock) GetIterator(terminate context.Context) stream.Iterator {
	//TODO implement me
	panic("implement me")
}

func (d *dataMock) Put(item communication.Item) {
	//TODO implement me
	panic("implement me")
}

func (d *dataMock) TryGet(position int) (error, communication.Item) {

	if d.counter == d.switchPosition {
		d.counter++
		return nil, communication.Item{Id: stream.NoItemId}
	}
	d.counter++
	return nil, communication.Item{Id: position}

}

func TestIteratorShouldSwitchToNotifier(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Testing iteration %d", i), testIteratorShouldSwitchToNotifier)
	}

}

func testIteratorShouldSwitchToNotifier(t *testing.T) {

	count := rand.Intn(100)

	t.Logf("Random is %d", count)

	notifier := make(chan int, count)

	for i := 0; i < count; i++ {
		notifier <- i
	}

	ctx := context.Background()

	dataSource := dataMock{
		switchPosition: 50,
	}

	iterator := stream.NewIterator(&dataSource, notifier, ctx)

	for j := 0; j < count; j++ {
		_, iteratorItem := iterator.Next()

		if (j < dataSource.switchPosition) && (iteratorItem.Source != "Iterator") {
			t.Errorf("Expected Iterator source recieved %s ", iteratorItem.Source)
		}
		if (j >= dataSource.switchPosition) && (iteratorItem.Source != "Notifier") {
			t.Errorf("Expected Notifier source recieved %s ", iteratorItem.Source)
		}

	}

}
