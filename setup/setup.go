package setup

import (
	"sync"
)

var Config PilsnerConfig

var initConfig sync.Once

func init() {
	initConfig.Do(func() {
		Config = PilsnerConfig{
			StoragePath:   "./",
			BufferSize:    Large,
			StorageOption: FileStorage,
		}
	})
}

type PilsnerConfig struct {

	// Defines path to directory with data
	StoragePath string

	// Size of in memory buffers
	BufferSize BufferSize

	// Store in memory
	StorageOption StorageType
}

type BufferSize int

const (
	Small BufferSize = 10
	Large            = 1000
)

type StorageType int

const (
	FileStorage   StorageType = 0
	MemoryStorage             = 1
)
