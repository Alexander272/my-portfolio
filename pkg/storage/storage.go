package storage

import (
	"context"
	"mime/multipart"
)

type Provider interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path, name string) (string, error)
}
