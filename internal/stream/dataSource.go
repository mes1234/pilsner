package stream

import "sync"

type Item struct {
	Id      int
	content []byte
}

type items struct {
	repository []Item
	lock       sync.Mutex
}

func NewDataSource() *items {
	repository := make([]Item, 0)

	return &items{repository: repository}
}

type Data interface {
	GetIterator() Iterator
	Put(item Item)
}

func (i *items) GetIterator() Iterator {
	return NewIterator(&i.repository)
}

func (i *items) Get(position int) Item {
	return i.repository[position]
}

func (i *items) Put(item Item) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.repository = append(i.repository, item)
}
