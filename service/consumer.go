package service

import (
	"fmt"
	"google.golang.org/grpc"
	"pilsner/handler"
	"pilsner/internal/communication"
	"pilsner/proto/pb"
	"pilsner/translator"
	"time"
)

type consumerService struct {
}

type consumerServiceHandler struct {
	h handler.ConsumerHandler
}

type Handler interface {
	Handle(server pb.Consumer_ConsumeServer) error
}

func (c *consumerService) Consume(server pb.Consumer_ConsumeServer) error {

	h := NewConsumerServiceHandler()

	return h.Handle(server)
}

func NewConsumerService() *consumerService {
	return &consumerService{}
}

func NewConsumerServiceHandler() *consumerServiceHandler {
	return &consumerServiceHandler{
		h: *handler.NewConsumerHandler(),
	}
}

func (c *consumerService) AttachTo(server *grpc.Server) {
	pb.RegisterConsumerServer(server, c)
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
