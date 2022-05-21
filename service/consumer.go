package service

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"pilsner/proto/pb"
)

type consumerService struct {
	streamingStartedFlag bool
}

func (c *consumerService) Consumer(server pb.Consumer_ConsumerServer) error {
	for {

		msg := pb.ConsumerResponse{}
		server.RecvMsg(&msg)

		log.Printf("Got msg")

		content := msg.GetContent()

		switch content.(type) {
		case *pb.ConsumerResponse_Setup:
			c.handleSetup(msg.GetSetup(), &server)
		case *pb.ConsumerResponse_Ack:
			c.handleAck(msg.GetAck(), &server)
		}
	}
}

func (c *consumerService) handleSetup(setup *pb.ConsumerSetup, server *pb.Consumer_ConsumerServer) {

}

func (c *consumerService) handleAck(setup *pb.ConsumerAck, server *pb.Consumer_ConsumerServer) {

}

func NewConsumerService() *consumerService {
	return &consumerService{
		streamingStartedFlag: false,
	}
}

func (c *consumerService) Publish(ctx context.Context, item *pb.Item) (*pb.ServerResponse, error) {

	response := pb.ServerResponse{}

	return &response, nil
}

func (c *consumerService) AttachTo(server *grpc.Server) {
	pb.RegisterConsumerServer(server, c)
}
