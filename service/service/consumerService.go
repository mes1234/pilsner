package service

import (
	"google.golang.org/grpc"
	"pilsner/proto/pb"
	"pilsner/service/handler"
)

type consumerService struct {
}

func (c *consumerService) Consume(server pb.Consumer_ConsumeServer) error {
	handler.NewConsumerServiceHandler().Handle(server)

	return nil
}

func (c *consumerService) AttachTo(server *grpc.Server) {
	pb.RegisterConsumerServer(server, c)
}

func NewConsumerService() *consumerService {
	return &consumerService{}
}
