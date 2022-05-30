package service

import (
	"google.golang.org/grpc"
	service "pilsner/handler"
	"pilsner/proto/pb"
)

type consumerService struct {
}

func (c *consumerService) Consume(server pb.Consumer_ConsumeServer) error {

	handler := service.NewConsumerServiceHandler()

	return handler.Handle(server)
}

func NewConsumerService() *consumerService {
	return &consumerService{}
}

func (c *consumerService) AttachTo(server *grpc.Server) {
	pb.RegisterConsumerServer(server, c)
}
