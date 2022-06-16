package handler

import (
	"context"
	"google.golang.org/grpc"
	"pilsner/proto/pb"
)

type Builder interface {
	AttachTo(server *grpc.Server)
}

type ConsumeServiceHandler interface {
	Handle(server pb.Consumer_ConsumeServer)
}

type PublishServiceHandler interface {
	Handle(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error)
}
