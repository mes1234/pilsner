package service

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"pilsner/internal/communication"
	"pilsner/internal/manager/consumerManager"
	"pilsner/proto/pb"
	"sync"
)

type consumerService struct {
	streamingStartedFlag bool
	Channel              <-chan communication.Item
	waiter               sync.WaitGroup
}

func (c *consumerService) Consume(server pb.Consumer_ConsumeServer) error {
	for {

		msg := pb.ConsumerResponse{}
		err := server.RecvMsg(&msg)
		if err != nil {
			return err
		}

		log.Printf("Got msg")

		content := msg.GetContent()

		switch content.(type) {
		case *pb.ConsumerResponse_Setup:
			c.handleSetup(msg.GetSetup())
		case *pb.ConsumerResponse_Ack:
			c.handleAck(msg.GetAck())
		default:
			return fmt.Errorf("not supported type")
		}

		go c.SendToConsumer(server)
	}
}

func (c *consumerService) SendToConsumer(server pb.Consumer_ConsumeServer) {

	if c.streamingStartedFlag != true {
		return
	}
	for {
		select {
		case data := <-c.Channel:
			p, ok := data.Content.([]byte)
			if ok {
				_ = server.Send(&pb.Item{Content: p})
			}
			c.waiter.Add(1)
			c.waiter.Wait()
		}
	}

}

func (c *consumerService) handleSetup(setup *pb.ConsumerSetup) error {

	if c.streamingStartedFlag == true {
		return fmt.Errorf("streaming already started")
	}

	manager := consumerManager.NewConsumerManager()

	_, delegate := manager.Attach(setup.StreamName, setup.ConsumerName)

	c.Channel = delegate.Channel

	c.streamingStartedFlag = true

	return nil
}

func (c *consumerService) handleAck(ack *pb.ConsumerAck) error {
	// TODO do status check
	// TODO this can lead to negative wait group counter
	c.waiter.Done()
	return nil
}

func NewConsumerService() *consumerService {
	return &consumerService{
		streamingStartedFlag: false,
	}
}

func (c *consumerService) AttachTo(server *grpc.Server) {
	pb.RegisterConsumerServer(server, c)
}
