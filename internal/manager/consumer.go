package manager

import (
	"pilsner/internal/communication"
	"sync"
)

var initConsumerManager sync.Once
var consumerManager ConsumerManager

func init() {
	initConsumerManager.Do(func() {
		consumerManager = dummyConsumerManager{
			streamManager: NewStreamManager(),
		}
	})
}

type dummyConsumerManager struct {
	streamManager StreamManager
}

func NewConsumerManager() ConsumerManager {
	return consumerManager
}

func (d dummyConsumerManager) Attach(streamName string, consumerName string) (error, communication.Delegate) {
	//TODO implement me
	panic("implement me")
}

type ConsumerManager interface {
	Attach(streamName string, consumerName string) (error, communication.Delegate)
}
