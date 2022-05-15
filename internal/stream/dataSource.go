package stream

import (
	"context"
	"github.com/google/uuid"
	"log"
	"pilsner/internal/communication"
	"sync"
)

const NoItemId = -1

type items struct {
	repository []communication.Item
	notifiers  map[uuid.UUID]chan int
	lock       sync.Mutex
	ctx        context.Context
}

func NewDataSource(ctx context.Context) *items {
	repository := make([]communication.Item, 0)

	return &items{
		repository: repository,
		ctx:        ctx,
		notifiers:  make(map[uuid.UUID]chan int),
	}
}

type Data interface {
	GetIterator(terminate context.Context) Iterator
	Put(item communication.Item)
	TryGet(position int) (error, communication.Item)
}

func (i *items) GetIterator(terminate context.Context) Iterator {

	notifier := i.addNotifier(terminate)

	return newIterator(i, notifier, terminate)
}

func (i *items) addNotifier(terminate context.Context) <-chan int {

	notifier := make(chan int, 2000)

	id := uuid.New()

	i.notifiers[id] = notifier

	// Schedule removal in case of termination
	go i.removeNotifier(id, terminate)

	return notifier
}

func (i *items) removeNotifier(id uuid.UUID, terminate context.Context) {
	select {
	case <-terminate.Done():
		i.lock.Lock()
		defer i.lock.Unlock()
		delete(i.notifiers, id)
		log.Printf("Notifier %s was removed", id)
	}
}

func (i *items) TryGet(position int) (error, communication.Item) {
	if len(i.repository)-1 >= position {
		return nil, i.repository[position]
	} else {
		return nil, communication.Item{
			Id: NoItemId,
		}
	}
}

func (i *items) Put(item communication.Item) {

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
