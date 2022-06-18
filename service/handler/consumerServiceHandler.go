package handler

import (
	"fmt"
	"pilsner/internal/adapter"
	"pilsner/internal/communication"
	"pilsner/internal/manager"
	"pilsner/proto/pb"
	"pilsner/translator"
)

type consumerServiceHandler struct {
	handler adapter.ConsumerHandler
}

func (c *consumerServiceHandler) Handle(server pb.Consumer_ConsumeServer) {

	adapter.ListenToConsumer[pb.ConsumerResponse](server.RecvMsg, c.buildHandleMsgFunction(&server))
}

func (c *consumerServiceHandler) buildHandleMsgFunction(server *pb.Consumer_ConsumeServer) adapter.HandleMsgFunction {

	return func(msg interface{}) error {

		val, ok := msg.(*pb.ConsumerResponse)

		if ok != true {
			return fmt.Errorf("wrong type expected Consumer Response")
		}

		content := val.GetContent()

		switch content.(type) {
		case *pb.ConsumerResponse_Setup:
			_, setupDto := translator.Translate[communication.ConsumerSetup](val.GetSetup())
			manager := manager.NewConsumerManager()
			_ = c.handler.HandleSetup(setupDto, manager)
			_, itemPb := translator.Translate[func(*communication.Item) error]((*server).Send)
			go c.handler.SendToConsumer(itemPb)

		case *pb.ConsumerResponse_Ack:
			_, ackDto := translator.Translate[communication.ConsumerAck](val.GetAck())
			_ = c.handler.HandleAck(ackDto)

		default:
			return fmt.Errorf("not supported type")
		}
		return nil
	}
}

func NewConsumerServiceHandler() ConsumeServiceHandler {
	return &consumerServiceHandler{
		handler: adapter.NewConsumerHandler(),
	}
}
