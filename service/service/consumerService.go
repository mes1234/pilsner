package service

import (
	"google.golang.org/grpc"
	"pilsner/proto/pb"
	"pilsner/service/handler"
)

type consumerService struct {
}

func (c *consumerService) Consume(server pb.Consumer_ConsumeServer) error {

	h := handler.NewConsumerServiceHandler()

	return h.Handle(server)
}

func (c *consumerService) AttachTo(server *grpc.Server) {
	pb.RegisterConsumerServer(server, c)
}

func NewConsumerService() *consumerService {
	return &consumerService{}
}
