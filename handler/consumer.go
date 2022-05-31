package handler

import (
	"context"
	"fmt"
	"log"
	"pilsner/internal/communication"
	"pilsner/mapper"
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

		err = c.handleMsg(&server, &msg)
		if err != nil {
			return
		}
	}
}

func (c *consumerServiceHandler) handleMsg(server *pb.Consumer_ConsumeServer, msg *pb.ConsumerResponse) error {
	content := msg.GetContent()

	switch content.(type) {
	case *pb.ConsumerResponse_Setup:
		_ = c.handleSetup(mapper.MapConsumerSetupProtoToInternal(msg.GetSetup()))
		go c.SendToConsumer(mapper.MapItemToProto((*server).Send))
	case *pb.ConsumerResponse_Ack:
		_ = c.handleAck(msg.GetAck())
	default:
		return fmt.Errorf("not supported type")
	}
	return nil
}

func (c *consumerServiceHandler) handleAck(ack *pb.ConsumerAck) error {
	log.Printf("Got Ack  %s", ack.Status.String())
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
