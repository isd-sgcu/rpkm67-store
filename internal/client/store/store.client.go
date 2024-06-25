package store

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type Client interface {
	Upload(bucketName, name string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	DeleteByKey(bucketName, key string, opts minio.RemoveObjectOptions) error
}

type clientImpl struct {
	*minio.Client
}

func (c *clientImpl) DeleteByKey(bucketName string, key string, opts minio.RemoveObjectOptions) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	return c.Client.RemoveObject(ctx, bucketName, key, opts)
}

func (c *clientImpl) Upload(bucketName string, name string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	return c.Client.PutObject(ctx, bucketName, name, reader, objectSize, opts)
}

func NewClient(minioClient *minio.Client) Client {
	return &clientImpl{minioClient}
}
