package streamManager

import (
	"context"
	"fmt"
	"pilsner/internal/stream"
	"sync"
)

var initStreamManager sync.Once
var streamManager StreamManager

func init() {
	initStreamManager.Do(func() {
		streamManager = &dummyStreamManager{
			streams: make(map[string]streamEntity),
		}
	})
}

type StreamManager interface {
	Add(name string) error
	Remove(name string) error
	Get(name string) (error, stream.StreamerPublisher)
}

type dummyStreamManager struct {
	lock    sync.RWMutex
	streams map[string]streamEntity
}

type streamEntity struct {
	stream stream.StreamerPublisher
	cancel context.CancelFunc
}

func NewStreamManager() StreamManager {
	return streamManager
}

func (d *dummyStreamManager) Add(name string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if _, ok := d.streams[name]; ok {
		return fmt.Errorf("stream %s already defined", name)
	}

	newStream, cancel := stream.NewStream()

	d.streams[name] = streamEntity{
		stream: newStream,
		cancel: cancel,
	}
	return nil

}

func (d *dummyStreamManager) Remove(name string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if _, ok := d.streams[name]; !ok {
		return fmt.Errorf("stream %s not defined", name)
	}

	delete(d.streams, name)
	return nil

}

func (d *dummyStreamManager) Get(name string) (error, stream.StreamerPublisher) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if instance, ok := d.streams[name]; ok {
		return nil, instance.stream
	} else {
		return fmt.Errorf("stream %s not defined", name), nil
	}
}
