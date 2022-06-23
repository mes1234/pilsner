package communication

import (
	"context"
	"pilsner/setup"
)

type Delegate struct {
	Context context.Context
	Channel chan Item
	Name    string
	Cancel  context.CancelFunc
}

func NewDelegate(name string) *Delegate {

	ctx, cancel := context.WithCancel(context.Background())

	subscriber := make(chan Item, setup.Config.BufferSize)

	return &Delegate{
		Name:    name,
		Channel: subscriber,
		Context: ctx,
		Cancel:  cancel,
	}
}
