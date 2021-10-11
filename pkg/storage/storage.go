package storage

import (
	"context"
	"mime/multipart"
)

type File struct {
	Name string
	Url  string
}

type Provider interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path, name string) (*File, error)
	Remove(ctx context.Context, path, filename string) error
}
