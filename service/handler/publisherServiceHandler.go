package handler

import (
	"context"
	"pilsner/internal/adapter"
	"pilsner/internal/communication"
	"pilsner/proto/pb"
	"pilsner/translator"
)

type publisherServiceHandler struct {
	handler adapter.PublisherHandler
}

func (p *publisherServiceHandler) Handle(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error) {

	_, itemDto := translator.Translate[communication.Item](item.Item)

	err := p.handler.Handle(itemDto, item.StreamName)

	if err != nil {
		return &pb.ServerResponse{
			Status: pb.AckStatusServer_Error,
		}, err
	} else {
		return &pb.ServerResponse{
			Status: pb.AckStatusServer_Received,
		}, nil
	}
}

func NewPublisherServiceHandler() PublishServiceHandler {
	return &publisherServiceHandler{
		handler: adapter.NewPublisherHandler(),
	}
}
