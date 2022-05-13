package stream

import (
	"context"
	"sync"
)

const NoItemId = -1

type Item struct {
	Id      int
	content []byte
	Source  string
}

type items struct {
	repository []Item
	notifiers  []chan int
	lock       sync.Mutex
	ctx        context.Context
}

func NewDataSource(ctx context.Context) *items {
	repository := make([]Item, 0)

	return &items{
		repository: repository,
		ctx:        ctx,
	}
}

type Data interface {
	GetIterator(terminate context.Context) Iterator
	Put(item Item)
	TryGet(position int) (error, Item)
}

func (i *items) GetIterator(terminate context.Context) Iterator {

	notifier := make(chan int, 2000)

	i.notifiers = append(i.notifiers, notifier)

	return NewIterator(i, notifier, terminate)
}

func (i *items) TryGet(position int) (error, Item) {
	if len(i.repository)-1 >= position {
		return nil, i.repository[position]
	} else {
		return nil, Item{
			Id: NoItemId,
		}
	}
}

func (i *items) Put(item Item) {

	select {
	case <-i.ctx.Done():
		return
	default:
		i.lock.Lock()
		defer i.lock.Unlock()
		i.repository = append(i.repository, item)

		for _, notifier := range i.notifiers {
			select {
			case notifier <- item.Id:
			}
		}
	}
}
