package main

import (
	"log"
	"net"
	"pilsner/internal/manager/consumerManager"
	"pilsner/internal/manager/streamManager"
	"pilsner/server"
	"pilsner/service"
)

func main() {

	// Init managers
	streamManager.NewStreamManager()
	consumerManager.NewConsumerManager()

	// Init some dummies stream
	sm := streamManager.NewStreamManager()
	sm.Add("Hello")

	log.Println("Starting listening on port 8080")
	port := ":8080"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)

	grpcServer := server.NewServer()

	publisherService := service.NewPublisherService()

	consumerService := service.NewConsumerService()

	publisherService.AttachTo(grpcServer)
	consumerService.AttachTo(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
