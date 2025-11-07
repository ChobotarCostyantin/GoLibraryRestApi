package storage

import (
	"encoding/json"
	"github.com/ChobotarCostyantin/GoLibraryRestApi/models"
	"os"
	"sync"
)

type Storage struct {
	Authors    []models.Author `json:"authors"`
	Books      []models.Book   `json:"books"`
	Librarians []models.Librarian `json:"librarians"`
	mutex      sync.RWMutex
}

var Store = &Storage{
	Authors:    []models.Author{},
	Books:      []models.Book{},
	Librarians: []models.Librarian{},
}

const storageFile = "library_data.json"

func (s *Storage) LoadFromFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := os.ReadFile(storageFile)
	if err != nil {
		return s.SaveToFile()
	}

	return json.Unmarshal(data, s)
}
func (s *Storage) SaveToFile() error {
	data, err := json.MarshalIndent(s, "", " ")

	if err != nil {
		return err
	}

	return os.WriteFile(storageFile, data, 0644)
}

func (s *Storage) Lock() {
	s.mutex.Lock()
}

func (s *Storage) Unlock() {
	s.mutex.Unlock()
}

func (s *Storage) RLock() {
	s.mutex.RLock()
}

func (s *Storage) RUnlock() {
	s.mutex.RUnlock()
}
