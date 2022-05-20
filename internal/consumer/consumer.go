package consumer

import (
	"pilsner/internal/communication"
	"time"
)

type Consumer interface {
}

type Callback func(item communication.Item) error
type policy func(item communication.Item, callback Callback) error

type consumer struct {
	callback Callback
	policy   policy
	delegate communication.Delegate
}

type Setup struct {
}

func defaultPolicy(item communication.Item, callback Callback) error {
	err := callback(item)
	return err
}

func NewConsumer(delegate communication.Delegate, callback Callback, setup Setup) *consumer {

	con := consumer{
		callback: callback,
		policy:   defaultPolicy,
		delegate: delegate,
	}

	go con.run()

	return &con
}

func (c *consumer) run() {
	for {
		select {
		case <-c.delegate.Context.Done():
			return
		case item := <-c.delegate.Channel:
			go c.policy(item, c.callback)
		case <-time.After(1 * time.Second):
		}
	}
}
