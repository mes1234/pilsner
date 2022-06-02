package handler

import (
	"fmt"
	"pilsner/internal/communication"
	"pilsner/internal/manager/consumerManager"
)

type SendFunction func(item *communication.Item) error

type ReceiveFunction func(m interface{}) error

type HandleMsgFunction func(msg interface{}) error

func (c *consumerServiceHandler) SendToConsumer(send SendFunction) {
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

func listenToConsumer[K interface{}](receive ReceiveFunction, handleMsg HandleMsgFunction) {
	for {
		var obj K
		err := receive(&obj)
		if err != nil {
			return
		}

		err = handleMsg(&obj)
		if err != nil {
			return
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
