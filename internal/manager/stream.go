package manager

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
	Get(name string) (error, stream.Publisher)
}

type dummyStreamManager struct {
	streams map[string]streamEntity
}

type streamEntity struct {
	stream stream.Publisher
	cancel context.CancelFunc
}

func NewStreamManager() StreamManager {
	return streamManager
}

func (d *dummyStreamManager) Add(name string) error {
	newStream, cancel := stream.NewStream()

	if _, ok := d.streams[name]; ok {
		return fmt.Errorf("stream %s already defined", name)
	}

	d.streams[name] = streamEntity{
		stream: newStream,
		cancel: cancel,
	}
	return nil

}

func (d *dummyStreamManager) Remove(name string) error {
	//TODO implement me
	panic("implement me")
}

func (d *dummyStreamManager) Get(name string) (error, stream.Publisher) {
	//TODO implement me
	panic("implement me")
}
