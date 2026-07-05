package storage

type MockStorage struct {
	storage map[int][]byte
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		storage: make(map[int][]byte),
	}
}

func (s *MockStorage) AddNewEntry(index int, hash []byte) {
	s.storage[index] = hash
}
