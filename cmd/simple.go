package main

import (
	"log"
	"net"
	"pilsner/internal/manager/consumerManager"
	"pilsner/internal/manager/streamManager"
	"pilsner/server"
	service2 "pilsner/service/service"
)

func main() {

	// Init managers
	streamManager.NewStreamManager()
	consumerManager.NewConsumerManager()

	// Init some dummies stream
	_ = streamManager.NewStreamManager()

	log.Println("Starting listening on port 8080")
	port := ":8080"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)

	grpcServer := server.NewServer()

	publisherService := service2.NewPublisherService()

	consumerService := service2.NewConsumerService()

	publisherService.AttachTo(grpcServer)
	consumerService.AttachTo(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
