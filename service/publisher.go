package service

import (
	"context"
	"google.golang.org/grpc"
	"pilsner/proto/pb"
)

type publisherService struct {
}

func NewPublisherService() *publisherService {
	return &publisherService{}
}

func (p *publisherService) Publish(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error) {

	response := pb.ServerResponse{}

	return &response, nil
}

func (p *publisherService) AttachTo(server *grpc.Server) {
	pb.RegisterPublisherServer(server, p)
}
