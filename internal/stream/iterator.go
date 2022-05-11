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

func (i *iterator) Next() (next bool, item Item) {
	for {
		if i.replayFinished {
			select {
			case <-i.terminator.Done():
				return false, Item{}
			case pos := <-i.notifier:
				if pos >= i.position {
					err, item := i.data.TryGet(pos)
					if err == nil {
						return true, Item{
							Id:      item.Id,
							content: item.content,
							Source:  "Notifier",
						}
					}
				}

			}
		} else {
			select {
			case <-i.terminator.Done():
				return false, Item{}
			default:
				err, item := i.data.TryGet(i.position)
				if err == nil {
					i.position++
					return true, Item{
						Id:      item.Id,
						content: item.content,
						Source:  "Iterator",
					}
				} else {
					i.replayFinished = true
				}
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
