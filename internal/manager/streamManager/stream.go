package streamManager

import (
	"context"
	"pilsner/internal/stream"
	"sync"
)

var initStreamManager sync.Once
var streamManager StreamManager

func init() {
	initStreamManager.Do(func() {
		newStream, cancel := stream.NewStream()

		streamManager = &dummyStreamManager{
			stream: streamEntity{
				stream: newStream,
				cancel: cancel,
			},
		}
	})
}

type StreamManager interface {
	Get() (error, stream.StreamerPublisher)
}

type dummyStreamManager struct {
	lock   sync.RWMutex
	stream streamEntity
}

type streamEntity struct {
	stream stream.StreamerPublisher
	cancel context.CancelFunc
}

func NewStreamManager() StreamManager {
	return streamManager
}

func (d *dummyStreamManager) Get() (error, stream.StreamerPublisher) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return nil, d.stream.stream

}
