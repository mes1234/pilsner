package service

import (
	"google.golang.org/grpc"
)

type Builder interface {
	AttachTo(server *grpc.Server)
}
