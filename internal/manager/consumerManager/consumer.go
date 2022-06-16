package consumerManager

import (
	"fmt"
	"pilsner/internal/communication"
	"pilsner/internal/manager/streamManager"
	"sync"
)

var initConsumerManager sync.Once
var consumerManager ConsumerManager

func init() {
	initConsumerManager.Do(func() {
		consumerManager = &dummyConsumerManager{
			streamManager: streamManager.NewStreamManager(),
			consumers:     make(map[string]communication.Delegate),
		}
	})
}

type dummyConsumerManager struct {
	lock          sync.RWMutex
	streamManager streamManager.StreamManager
	consumers     map[string]communication.Delegate
}

func NewConsumerManager() ConsumerManager {
	return consumerManager
}

func (d *dummyConsumerManager) Attach(streamName string, consumerName string) (error, communication.Delegate) {
	d.lock.Lock()
	defer d.lock.Unlock()

	err, s := d.streamManager.Get(streamName)

	if err != nil {
		return fmt.Errorf("stream %s not defined", streamName), communication.Delegate{}
	}

	if _, ok := d.consumers[consumerName]; ok {
		return fmt.Errorf("consumer %s already defined", consumerName), communication.Delegate{}
	}

	delegate := communication.NewDelegate(consumerName)

	s.Stream(*delegate)

	d.consumers[consumerName] = *delegate

	return nil, *delegate

}

type ConsumerManager interface {
	Attach(streamName string, consumerName string) (error, communication.Delegate)
}
