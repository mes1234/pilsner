package stream

import (
	"log"
	"os"
	"pilsner/internal/communication"
	"sync"
)

const (
	indexFileName = "index.pilsner"
	dataFileName  = "data.pilsner"
)

var indexFile *os.File
var dataFile *os.File

var instance fileRepository
var initFileRepository sync.Once

type fileRepository struct {
	path string
}

func NewFileRepository(path string) *fileRepository {

	initFileRepository.Do(func() {
		indexFile = OpenOrCreate(path, indexFileName)
		dataFile = OpenOrCreate(path, dataFileName)
	})

	return &instance
}

func (m *fileRepository) Append(item communication.Item) RepositoryWrapper {

	return m
}

func (m *fileRepository) Len() int {
	return 0
}

func (m *fileRepository) Get(position int) communication.Item {
	return communication.Item{}
}

func OpenOrCreate(path string, fileName string) *os.File {

	file, err := os.OpenFile(path+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("Cannot open file %s", path+fileName)
	}

	return file
}
