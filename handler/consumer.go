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

	go listenToConsumer[pb.ConsumerResponse](server.RecvMsg, c.buildHandleMsgFunction(&server))

	for {
		select {
		case <-c.ctx.Done():
			return fmt.Errorf("closed consumer stream")
		case <-time.After(1 * time.Second):

		}
	}
}
func (c *consumerServiceHandler) buildHandleMsgFunction(server *pb.Consumer_ConsumeServer) HandleMsgFunction {

	return func(msg interface{}) error {

		val, ok := msg.(*pb.ConsumerResponse)

		if ok != true {
			return fmt.Errorf("wrong type expected Consumer Response")
		}

		content := val.GetContent()

		switch content.(type) {
		case *pb.ConsumerResponse_Setup:
			_ = c.handleSetup(mapper.MapConsumerSetupProtoToInternal(val.GetSetup()))
			go c.SendToConsumer(mapper.MapItemToProto((*server).Send))
		case *pb.ConsumerResponse_Ack:
			_ = c.handleAck(val.GetAck())
		default:
			return fmt.Errorf("not supported type")
		}
		return nil
	}
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
