package handler

import (
	"context"
	"pilsner/internal/communication"
	"pilsner/internal/handler"
	"pilsner/proto/pb"
	"pilsner/translator"
)

type publisherServiceHandler struct {
	h handler.PublisherHandler
}

func (p *publisherServiceHandler) Handle(ctx context.Context, item *pb.PublisherRequest) (*pb.ServerResponse, error) {

	_, itemDto := translator.Translate[communication.Item](item.Item)

	err := p.h.Handle(itemDto, item.StreamName)

	return &pb.ServerResponse{}, err

}

func NewPublisherServiceHandler() *publisherServiceHandler {
	return &publisherServiceHandler{
		h: *handler.NewPublisherHandler(),
	}
}
