package subscriber

import (
	"github.com/google/uuid"
	"pilsner/internal/stream"
)

/////////////
// PUBLIC  //
/////////////

type Callback func(item stream.Item) error

func NewSubscriber(callback Callback, iterator stream.Iterator) *subscriber {
	sub := subscriber{
		id:       uuid.New(),
		callback: callback,
		iterator: iterator,
	}

	go sub.run()

	return &sub
}

/////////////
// PRIVATE //
/////////////

type subscriber struct {
	callback Callback
	id       uuid.UUID
	iterator stream.Iterator
}

// start
func (s *subscriber) run() {
	channel := make(chan stream.Item)

	s.iterator.Start(channel)

	for elem := range channel {
		s.callback(elem)
	}

}
