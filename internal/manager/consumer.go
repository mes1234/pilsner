package manager

import (
	"fmt"
	"pilsner/internal/communication"
	"pilsner/internal/stream"
	"sync"
)

var initConsumerManager sync.Once
var consumerManager ConsumerManager

func init() {
	initConsumerManager.Do(func() {
		consumerManager = &dummyConsumerManager{
			stream:    stream.Get(),
			consumers: make(map[string]communication.Delegate),
		}
	})
}

type dummyConsumerManager struct {
	lock      sync.RWMutex
	stream    stream.StreamerPublisher
	consumers map[string]communication.Delegate
}

func NewConsumerManager() ConsumerManager {
	return consumerManager
}

func (d *dummyConsumerManager) Attach(consumerName string) (error, communication.Delegate) {
	d.lock.Lock()
	defer d.lock.Unlock()

	s := stream.Get()

	if _, ok := d.consumers[consumerName]; ok {
		return fmt.Errorf("consumer %s already defined", consumerName), communication.Delegate{}
	}

	delegate := communication.NewDelegate(consumerName)

	s.Stream(*delegate)

	d.consumers[consumerName] = *delegate

	return nil, *delegate

}

type ConsumerManager interface {
	Attach(consumerName string) (error, communication.Delegate)
}
