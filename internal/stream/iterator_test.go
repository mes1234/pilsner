package stream_test

import (
	"context"
	"pilsner/internal/communication"
	"pilsner/internal/stream"
	"testing"
)

type dataMock struct {
	counter int
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

	if d.counter == 1 {
		d.counter++
		return nil, communication.Item{Id: stream.NoItemId}
	}
	d.counter++
	return nil, communication.Item{Id: position}

}

// TODO fuzz it
func TestIteratorShouldSwitchToNotifier(t *testing.T) {

	notifier := make(chan int, 10)

	notifier <- 1
	notifier <- 2

	ctx := context.Background()

	dataSource := dataMock{}

	iterator := stream.NewIterator(&dataSource, notifier, ctx)

	_, iteratorItem := iterator.Next()

	_, notifierItem1 := iterator.Next()

	_, notifierItem2 := iterator.Next()

	if iteratorItem.Source != "Iterator" {
		t.Errorf("Expected Iterator source recieved %s ", iteratorItem.Source)
	}

	if notifierItem1.Source != "Notifier" {
		t.Errorf("Expected Notifier source recieved %s ", notifierItem1.Source)
	}
	if notifierItem2.Source != "Notifier" {
		t.Errorf("Expected Notifier source recieved %s ", notifierItem1.Source)
	}
}
