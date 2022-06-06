package handler

import (
	"context"
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/manager/streamManager"
	"pilsner/proto/pb"
)

type publisherServiceHandler struct {
}

func (p *publisherServiceHandler) Handle(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error) {

	stream := streamManager.NewStreamManager()

	_, publisher := stream.Get(item.StreamName)

	_ = publisher.Publish(communication.Item{
		Content: item.GetItem().Content,
	})

	log.Printf("Published to stream %s item", item.StreamName)

	return &pb.ServerResponse{}, nil

}

func NewPublisherServiceHandler() *publisherServiceHandler {
	return &publisherServiceHandler{}
}
