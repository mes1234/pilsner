package service

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/manager"
	"pilsner/proto/pb"
)

type consumerService struct {
	streamingStartedFlag bool
	Channel              <-chan communication.Item
}

func (c *consumerService) Consume(server pb.Consumer_ConsumeServer) error {
	for {

		msg := pb.ConsumerResponse{}
		err := server.RecvMsg(&msg)
		if err != nil {
			return err
		}

		log.Printf("Got msg")

		content := msg.GetContent()

		switch content.(type) {
		case *pb.ConsumerResponse_Setup:
			c.handleSetup(msg.GetSetup())
		case *pb.ConsumerResponse_Ack:
			c.handleAck(msg.GetAck(), server)
		default:
			return fmt.Errorf("not supported type")
		}
	}
}

func (c *consumerService) handleSetup(setup *pb.ConsumerSetup) error {

	manager := manager.NewConsumerManager()

	_, delegate := manager.Attach(setup.StreamName, setup.ConsumerName)

	c.Channel = delegate.Channel

	return nil
}

func (c *consumerService) handleAck(setup *pb.ConsumerAck, server pb.Consumer_ConsumeServer) error {
	for {
		select {
		case data := <-c.Channel:
			p, ok := data.Content.([]byte)
			if ok {
				_ = server.Send(&pb.Item{Content: p})
			}
		}
	}
}

func NewConsumerService() *consumerService {
	return &consumerService{
		streamingStartedFlag: false,
	}
}

func (c *consumerService) AttachTo(server *grpc.Server) {
	pb.RegisterConsumerServer(server, c)
}
