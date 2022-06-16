package adapter

import (
	"context"
	"fmt"
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/manager/consumerManager"
)

type consumerHandler struct {
	startedFlag  bool
	finishedFlag bool
	Channel      <-chan communication.Item
	Waiter       chan struct{}
	Ctx          context.Context
}

type ConsumerHandler interface {
	SendToConsumer(send SendFunction)
	HandleSetup(setup communication.ConsumerSetup, manager consumerManager.ConsumerManager) error
	HandleAck(ack communication.ConsumerAck) error
}

func NewConsumerHandler() *consumerHandler {
	return &consumerHandler{
		startedFlag: false,
		Ctx:         context.Background(),
		Waiter:      make(chan struct{}),
	}
}

type SendFunction func(item *communication.Item) error

type ReceiveFunction func(m interface{}) error

type HandleMsgFunction func(msg interface{}) error

func (c *consumerHandler) SendToConsumer(send SendFunction) {

	if c.startedFlag != true {
		return
	}
	for {
		select {
		case <-c.Ctx.Done():
			c.finishedFlag = true
			return
		case data := <-c.Channel:
			p, ok := data.Content.([]byte)
			if ok {
				_ = send(&communication.Item{Content: p})
			}
			<-c.Waiter
		}
	}
}

func ListenToConsumer[K interface{}](receive ReceiveFunction, handleMsg HandleMsgFunction) {
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

func (c *consumerHandler) HandleSetup(setup communication.ConsumerSetup, manager consumerManager.ConsumerManager) error {

	if c.startedFlag == true {
		return fmt.Errorf("streaming already started")
	}

	_, delegate := manager.Attach(setup.StreamName, setup.ConsumerName)

	c.Channel = delegate.Channel
	c.Ctx = delegate.Context

	c.startedFlag = true

	return nil
}

func (c *consumerHandler) HandleAck(ack communication.ConsumerAck) error {
	log.Printf("Got Ack  %s", ack.String())
	c.Waiter <- struct{}{}
	return nil
}
