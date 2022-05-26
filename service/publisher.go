package service

import (
	"context"
	"google.golang.org/grpc"
	"pilsner/internal/communication"
	"pilsner/internal/manager/streamManager"
	"pilsner/proto/pb"
)

type publisherService struct {
}

func NewPublisherService() *publisherService {
	return &publisherService{}
}

func (p *publisherService) Publish(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error) {

	_ = p.handlePublisherRequest(item)

	return &pb.ServerResponse{}, nil

}

func (p *publisherService) handlePublisherRequest(item *pb.PublisherRequest) error {

	stream := streamManager.NewStreamManager()

	_, publisher := stream.Get(item.StreamName)

	_ = publisher.Publish(communication.Item{})

	return nil

}

func (p *publisherService) AttachTo(server *grpc.Server) {
	pb.RegisterPublisherServer(server, p)
}
