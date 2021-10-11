package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/Alexander272/my-portfolio/pkg/logger"
	"github.com/chai2010/webp"
	"google.golang.org/api/option"
)

type FileStorage struct {
	storage    *storage.BucketHandle
	bucketName string
}

func NewFileStorage(bucketName, pathToCredentials string) (*FileStorage, error) {
	config := &firebase.Config{
		StorageBucket: bucketName + ".appspot.com",
	}
	opt := option.WithCredentialsFile(pathToCredentials)
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		storage:    bucket,
		bucketName: bucketName,
	}, nil
}

func (fs *FileStorage) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path, name string) (string, error) {
	fileBytes, err := io.ReadAll(file)
	logger.Debug(header.Header)
	if err != nil {
		return "", err
	}

	newFile, err := fs.imageCompressing(fileBytes, 85, header.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	var filename string
	if name != "" {
		filename = name + ".webp"
	} else {
		filename = strings.Split(header.Filename, ".")[0] + fmt.Sprintf("_%d.webp", time.Now().Unix())
	}

	wc := fs.storage.Object(filepath.Join(path, filename)).NewWriter(ctx)
	_, err = io.Copy(wc, bytes.NewReader(newFile))
	if err != nil {
		return "", err

	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://storage.cloud.google.com/%s.appspot.com/%s", fs.bucketName, filepath.Join(path, filename)), nil
}

func (fs *FileStorage) imageCompressing(buffer []byte, quality float32, contentType string) ([]byte, error) {
	var img image.Image
	var err error
	switch contentType {
	case "image/png":
		img, err = png.Decode(bytes.NewReader(buffer))
		if err != nil {
			return nil, err
		}
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(bytes.NewReader(buffer))
		if err != nil {
			return nil, err
		}
	}

	var out bytes.Buffer

	if err = webp.Encode(&out, img, &webp.Options{Lossless: true, Exact: true, Quality: quality}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil

	// converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	// if err != nil {
	// 	return nil, err
	// }

	// processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	// if err != nil {
	// 	return nil, err
	// }

	// return processed, nil
}
