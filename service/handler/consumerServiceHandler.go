package handler

import (
	"fmt"
	"pilsner/internal/communication"
	"pilsner/internal/handler"
	"pilsner/proto/pb"
	"pilsner/translator"
	"time"
)

type consumerServiceHandler struct {
	h handler.ConsumerHandler
}

func (c *consumerServiceHandler) Handle(server pb.Consumer_ConsumeServer) error {

	go handler.ListenToConsumer[pb.ConsumerResponse](server.RecvMsg, c.buildHandleMsgFunction(&server))

	for {
		select {
		case <-c.h.Ctx.Done():
			return fmt.Errorf("closed consumer stream")
		case <-time.After(1 * time.Second):
		}
	}
}

func (c *consumerServiceHandler) buildHandleMsgFunction(server *pb.Consumer_ConsumeServer) handler.HandleMsgFunction {

	return func(msg interface{}) error {

		val, ok := msg.(*pb.ConsumerResponse)

		if ok != true {
			return fmt.Errorf("wrong type expected Consumer Response")
		}

		content := val.GetContent()

		switch content.(type) {
		case *pb.ConsumerResponse_Setup:
			_, setupDto := translator.Translate[communication.ConsumerSetup](val.GetSetup())
			_ = c.h.HandleSetup(setupDto)

			_, itemPb := translator.Translate[func(*communication.Item) error]((*server).Send)
			go c.h.SendToConsumer(itemPb)

		case *pb.ConsumerResponse_Ack:
			_, ackDto := translator.Translate[communication.ConsumerAck](val.GetAck())

			_ = c.h.HandleAck(ackDto)
		default:
			return fmt.Errorf("not supported type")
		}
		return nil
	}
}

func NewConsumerServiceHandler() *consumerServiceHandler {
	return &consumerServiceHandler{
		h: *handler.NewConsumerHandler(),
	}
}
