package service

import (
	"context"
	"google.golang.org/grpc"
	"pilsner/proto/pb"
)

type Builder interface {
	AttachTo(server *grpc.Server)
}

type ConsumeHandler interface {
	Handle(server pb.Consumer_ConsumeServer) error
}

type PublishHandler interface {
	Handle(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error)
}
