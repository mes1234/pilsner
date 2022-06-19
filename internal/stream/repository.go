package stream

import "pilsner/internal/communication"

type RepositoryWrapper interface {
	Append(item communication.Item) RepositoryWrapper
	Len() int
	Get(position int) communication.Item
}

type memoryRepository struct {
	repository []communication.Item
}

func NewRepository() *memoryRepository {
	return &memoryRepository{repository: make([]communication.Item, 0)}
}

func (m *memoryRepository) Append(item communication.Item) RepositoryWrapper {
	m.repository = append(m.repository, item)
	return m
}

func (m *memoryRepository) Len() int {
	return len(m.repository)
}

func (m *memoryRepository) Get(position int) communication.Item {
	return m.repository[position]
}
