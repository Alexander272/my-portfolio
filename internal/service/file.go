package service

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/Alexander272/my-portfolio/internal/domain"
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

func (s *FileService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path, filename string) (*domain.File, error) {
	if !strings.Contains(header.Header.Get("Content-Type"), "image") {
		return nil, nil
	}
	res, err := s.storage.Upload(ctx, file, header, path, filename)
	if err != nil {
		return nil, err
	}
	return &domain.File{
		FileType: "Image",
		Name:     res.Name,
		OrigName: header.Filename,
		Url:      res.Url,
	}, nil
}

func (s *FileService) Remove(ctx context.Context, path, filename string) error {
	return s.storage.Remove(ctx, path, filename)
}
