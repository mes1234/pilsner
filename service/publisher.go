package service

import (
	"context"
	"google.golang.org/grpc"
	"pilsner/proto/pb"
)

type publisherService struct {
}

func (p *publisherService) Publish(ctx context.Context, item *pb.Item) (*pb.ServerResponse, error) {

	response := pb.ServerResponse{}

	return &response, nil
}

func NewPublisherService() *grpc.Server {
	grpcServer := grpc.NewServer()

	service := publisherService{}
	pb.RegisterPublisherServer(grpcServer, &service)

	return grpcServer
}
