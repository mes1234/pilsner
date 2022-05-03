package consumer

import (
	"pilsner/internal/stream"
)

const DefaultRetryAttempts int = 3

type memoryConsumer struct {
	stream        <-chan stream.Item
	ConsumedItems int
	FailedItems   int
	callback      Callback
	retryPolicy   RetryPolicy
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

func NewConsumer(stream <-chan stream.Item, callback Callback) *memoryConsumer {
	consumer := memoryConsumer{
		stream:      stream,
		retryPolicy: simpleNRetryPolicy(DefaultRetryAttempts),
		callback:    callback,
	}

	go consumer.startProcessing()

	return &consumer
}
