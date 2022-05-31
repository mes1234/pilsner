package handler

import (
	"fmt"
	"pilsner/internal/communication"
	"pilsner/internal/manager/consumerManager"
)

func (c *consumerServiceHandler) SendToConsumer(send func(item *communication.Item) error) {
	if c.startedFlag != true {
		return
	}
	for {
		select {
		case <-c.ctx.Done():
			c.finishedFlag = true
			return
		case data := <-c.Channel:
			p, ok := data.Content.([]byte)
			if ok {
				_ = send(&communication.Item{Content: p})
			}
			<-c.waiter
		}
	}
}

func (c *consumerServiceHandler) handleSetup(setup communication.ConsumerSetup) error {

	if c.startedFlag == true {
		return fmt.Errorf("streaming already started")
	}

	manager := consumerManager.NewConsumerManager()

	_, delegate := manager.Attach(setup.StreamName, setup.ConsumerName)

	c.Channel = delegate.Channel
	c.ctx = delegate.Context

	c.startedFlag = true

	return nil
}
