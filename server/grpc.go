package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

type grpcServer struct {
}

type Server interface {
	Attach()
}

func DummyAuth(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func NewServer() *grpc.Server {

	unaryInterceptor := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpcAuth.UnaryServerInterceptor(DummyAuth)))

	streamInterceptor := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(grpcAuth.StreamServerInterceptor(DummyAuth)))

	opts := []grpc.ServerOption{
		unaryInterceptor,
		streamInterceptor,
	}
	grpcServer := grpc.NewServer(opts...)

	return grpcServer
}
