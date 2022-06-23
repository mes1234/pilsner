package stream

import (
	"pilsner/internal/communication"
	"pilsner/setup"
)

type RepositoryWrapper interface {
	Append(item communication.Item) RepositoryWrapper
	Len() int
	Get(position int) communication.Item
}

func NewRepository() RepositoryWrapper {
	storageOption := setup.Config.StorageOption

	switch storageOption {
	case setup.FileStorage:
		return NewFileRepository(setup.Config.StoragePath)
	case setup.MemoryStorage:
		return NewMemoryRepository()
	default:
		panic("Not supported storage type")
	}
}
