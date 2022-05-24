package manager

import (
	"pilsner/internal/stream"
	"sync"
)

var initStreamManager sync.Once
var streamManager StreamManager

func init() {
	initStreamManager.Do(func() {
		streamManager = dummyStreamManager{}
	})
}

type StreamManager interface {
	Add(name string) error
	Remove(name string) error
	Get(name string) (error, stream.Publisher)
}

type dummyStreamManager struct {
}

func NewStreamManager() StreamManager {
	return streamManager
}

func (d dummyStreamManager) Add(name string) error {
	//TODO implement me
	panic("implement me")
}

func (d dummyStreamManager) Remove(name string) error {
	//TODO implement me
	panic("implement me")
}

func (d dummyStreamManager) Get(name string) (error, stream.Publisher) {
	//TODO implement me
	panic("implement me")
}
