package service

import (
	"context"
	"google.golang.org/grpc"
	"pilsner/proto/pb"
	"pilsner/service/handler"
)

type publisherService struct {
}

func NewPublisherService() *publisherService {
	return &publisherService{}
}

func (p *publisherService) Publish(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error) {

	h := handler.NewPublisherServiceHandler()

	return h.Handle(ctx, item)

}

func (p *publisherService) AttachTo(server *grpc.Server) {
	pb.RegisterPublisherServer(server, p)
}
