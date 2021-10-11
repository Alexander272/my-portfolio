package service

import (
	"context"
	"mime/multipart"

	"github.com/Alexander272/my-portfolio/pkg/storage"
)

type FileService struct {
	storage storage.Provider
}

func NewFileService(storage storage.Provider) *FileService {
	return &FileService{
		storage: storage,
	}
}

func (s *FileService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	return s.storage.Upload(ctx, file, header, "avatar", "avatar")
}
