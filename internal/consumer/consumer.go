package consumer

import (
	"pilsner/internal/stream"
)

const DefaultRetryAttempts int = 3

type Consumer interface {
	RegisterCallback(callback Callback)
}

type memoryConsumer struct {
	stream        <-chan stream.Item
	ConsumedItems int
	FailedItems   int
	callback      Callback
	retryPolicy   RetryPolicy
}

func (c *memoryConsumer) RegisterCallback(callback Callback) {
	c.callback = callback
}

func (c *memoryConsumer) startProcessing() {
	for item := range c.stream {
		if err := c.retryPolicy(c.callback, item); err == nil {
			c.ConsumedItems++
		} else {
			c.FailedItems++
		}
	}
}

func simpleNRetryPolicy(n int) RetryPolicy {
	return func(callback Callback, item stream.Item) (err error) {
		for attempt := 0; attempt < n; attempt++ {
			if err = callback(item); err == nil {
				break
			}
		}
		return err
	}
}

func breakerCallback(item stream.Item) error {
	return nil
}

func NewConsumer(stream <-chan stream.Item) *memoryConsumer {
	consumer := memoryConsumer{
		stream:      stream,
		retryPolicy: simpleNRetryPolicy(DefaultRetryAttempts),
		callback:    breakerCallback,
	}

	go consumer.startProcessing()

	return &consumer
}
