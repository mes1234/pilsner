package communication

import (
	"pilsner/internal/setup"
)

type Delegate struct {
	Channel chan Item
	Name    string
}

func NewDelegate(name string) *Delegate {

	subscriber := make(chan Item, setup.BufferSize)

	return &Delegate{
		Name:    name,
		Channel: subscriber,
	}
}
