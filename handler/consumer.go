package service

import (
	"context"
	"fmt"
	"pilsner/internal/communication"
	"pilsner/internal/manager/consumerManager"
	"pilsner/proto/pb"
	"time"
)

type Handler interface {
	Handle(server pb.Consumer_ConsumeServer) error
}

type consumerServiceHandler struct {
	startedFlag  bool
	finishedFlag bool
	Channel      <-chan communication.Item
	waiter       chan struct{}
	ctx          context.Context
}

func (c *consumerServiceHandler) Handle(server pb.Consumer_ConsumeServer) error {

	go c.listenToConsumer(server)

	for {
		select {
		case <-c.ctx.Done():
			return fmt.Errorf("closed consumer stream")
		case <-time.After(1 * time.Second):

		}
	}
}

func (c *consumerServiceHandler) listenToConsumer(server pb.Consumer_ConsumeServer) {
	for {
		msg := pb.ConsumerResponse{}
		err := server.RecvMsg(&msg)
		if err != nil {
			return
		}

		err = c.handleMsg(server, msg)
		if err != nil {
			return
		}
	}
}

func (c *consumerServiceHandler) handleMsg(server pb.Consumer_ConsumeServer, msg pb.ConsumerResponse) error {
	content := msg.GetContent()

	switch content.(type) {
	case *pb.ConsumerResponse_Setup:
		c.handleSetup(msg.GetSetup())
		go c.SendToConsumer(server)
	case *pb.ConsumerResponse_Ack:
		c.handleAck(msg.GetAck())
	default:
		return fmt.Errorf("not supported type")
	}
	return nil
}

func (c *consumerServiceHandler) SendToConsumer(server pb.Consumer_ConsumeServer) {
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
				_ = server.Send(&pb.Item{Content: p})
			}
			<-c.waiter
		}
	}
}

func (c *consumerServiceHandler) handleSetup(setup *pb.ConsumerSetup) error {

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

func (c *consumerServiceHandler) handleAck(ack *pb.ConsumerAck) error {
	// TODO do status check
	// TODO this can lead to negative wait group counter
	c.waiter <- struct{}{}
	return nil
}

func NewConsumerServiceHandler() *consumerServiceHandler {
	return &consumerServiceHandler{
		startedFlag: false,
		ctx:         context.Background(),
		waiter:      make(chan struct{}),
	}
}
