package service

import "github.com/Alexander272/my-portfolio/pkg/storage"

type FileService struct {
	storage storage.Provider
}

func NewFileService(storage storage.Provider) *FileService {
	return &FileService{
		storage: storage,
	}
}
