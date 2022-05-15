package stream

import (
	"context"
	"pilsner/internal/communication"
)

type Iterator interface {
	Next() (next bool, item communication.Item)
}

type iterator struct {
	position       int
	data           Data
	terminator     context.Context
	notifier       <-chan int
	replayFinished bool
}

func (i *iterator) observeStream() (skip bool, next bool, item communication.Item) {
	select {
	case pos := <-i.notifier:

		if pos < i.position {
			// Discard as old and already processed
			return true, true, communication.Item{}
		}

		err, item := i.data.TryGet(pos)

		if err != nil {
			// Error occurred and consumer.go is failed
			return false, false, communication.Item{}
		}

		if item.Id == NoItemId {
			i.replayFinished = true
			return true, false, communication.Item{}
		}

		// Successful data
		return false, true, communication.Item{
			Id:      item.Id,
			Content: item.Content,
			Source:  "Notifier",
		}
	}
}

func (i *iterator) replayStream() (skip bool, next bool, item communication.Item) {

	err, item := i.data.TryGet(i.position)

	if err != nil {
		// Error occurred and consumer.go is failed
		return false, false, communication.Item{}
	}

	if item.Id == NoItemId {
		i.replayFinished = true
		return true, false, communication.Item{}
	}

	i.position++
	return false, true, communication.Item{
		Id:      item.Id,
		Content: item.Content,
		Source:  "Iterator",
	}

}

func (i *iterator) Next() (next bool, item communication.Item) {
	for {
		select {
		case <-i.terminator.Done():
			// Terminate processing
			return false, communication.Item{}
		default:
			var skip, next bool
			var item communication.Item

			if i.replayFinished {
				skip, next, item = i.observeStream()
			} else {
				skip, next, item = i.replayStream()
			}
			if !skip {
				return next, item
			}
		}
	}

}

func newIterator(data Data, notifier <-chan int, terminate context.Context) Iterator {
	return &iterator{
		data:           data,
		terminator:     terminate,
		notifier:       notifier,
		replayFinished: false,
	}
}
