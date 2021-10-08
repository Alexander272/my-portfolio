package storage

import (
	"context"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
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
