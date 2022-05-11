package stream

import (
	"context"
	"fmt"
	"sync"
)

type Item struct {
	Id      int
	content []byte
	Source  string
}

type items struct {
	repository []Item
	notifiers  []chan int
	lock       sync.Mutex
}

func NewDataSource() *items {
	repository := make([]Item, 0)

	return &items{repository: repository}
}

type Data interface {
	GetIterator(terminate context.Context) Iterator
	Put(item Item)
	TryGet(position int) (error, Item)
}

func (i *items) GetIterator(terminate context.Context) Iterator {

	notifier := make(chan int, 10)

	i.notifiers = append(i.notifiers, notifier)

	return NewIterator(i, notifier, terminate)
}

func (i *items) TryGet(position int) (error, Item) {
	if len(i.repository)-1 >= position {
		return nil, i.repository[position]
	} else {
		return fmt.Errorf("no Item"), Item{}
	}
}

func (i *items) Put(item Item) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.repository = append(i.repository, item)

	for _, notifier := range i.notifiers {
		notifier <- item.Id
	}
}
