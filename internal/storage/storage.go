package storage

import (
	"encoding/json"
	"os"
	"study/internal/models"
)

type Storage interface {
	SaveItems(items []models.Item) error
}

type FilesStorage struct {
	filename string
}

func NewFilesStorage(filename string) *FilesStorage {
	return &FilesStorage{filename: filename}
}

func (f *FilesStorage) SaveItems(items []models.Item) error {
	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.filename, data, 0644)
}

func (f *FilesStorage) LoadItems() ([]models.Item, error) {
	data, err := os.ReadFile(f.filename)
	if err != nil {
		return nil, err
	}

	var items []models.Item
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, err
}
