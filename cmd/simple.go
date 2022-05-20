package main

import (
	"log"
	"net"
	"pilsner/service"
)

func main() {

	log.Println("Starting listening on port 8080")
	port := ":8080"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)

	grpcServerPublisher := service.NewPublisherService()

	if err := grpcServerPublisher.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
