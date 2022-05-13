package stream

import (
	"context"
)

type Iterator interface {
	Next() (next bool, item Item)
}

type iterator struct {
	position       int
	data           Data
	terminator     context.Context
	notifier       <-chan int
	replayFinished bool
}

func (i *iterator) observeNotifier() (skip bool, next bool, item Item) {
	select {
	case pos := <-i.notifier:

		if pos < i.position {
			// Discard as old and already processed
			return true, true, Item{}
		}

		err, item := i.data.TryGet(pos)

		if err != nil {
			// Error occurred and consumer is failed
			return false, false, Item{}
		}

		// Successful data
		return false, true, Item{
			Id:      item.Id,
			content: item.content,
			Source:  "Notifier",
		}
	}
}

func (i *iterator) replayStream() (skip bool, next bool, item Item) {

	err, item := i.data.TryGet(i.position)

	if err != nil {
		// Error occurred and consumer is failed
		return false, false, Item{}
	}

	if item.Id == NoItemId {
		i.replayFinished = true
		return true, false, Item{}
	}

	i.position++
	return false, true, Item{
		Id:      item.Id,
		content: item.content,
		Source:  "Iterator",
	}

}

func (i *iterator) Next() (next bool, item Item) {
	for {
		select {
		case <-i.terminator.Done():
			// Terminate processing
			return false, Item{}
		default:
			var skip, next bool
			var item Item

			if i.replayFinished {
				skip, next, item = i.observeNotifier()
			} else {
				skip, next, item = i.replayStream()
			}
			if !skip {
				return next, item
			}
		}
	}

}

func NewIterator(data Data, notifier <-chan int, terminate context.Context) Iterator {
	return &iterator{
		data:           data,
		terminator:     terminate,
		notifier:       notifier,
		replayFinished: false,
	}
}
