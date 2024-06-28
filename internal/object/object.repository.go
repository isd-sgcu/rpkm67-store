package object

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/isd-sgcu/rpkm67-store/config"
	httpClient "github.com/isd-sgcu/rpkm67-store/internal/client/http"
	storeClient "github.com/isd-sgcu/rpkm67-store/internal/client/store"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

type Repository interface {
	Upload(file []byte, bucketName string, objectKey string) (url string, key string, err error)
	Delete(bucketName string, objectKey string) (err error)
	Get(bucketName string, objectKey string) (url string, err error)
	GetURL(bucketName string, objectKey string) string
}

type repositoryImpl struct {
	conf        *config.Store
	storeClient storeClient.Client
	httpClient  httpClient.Client
}

func NewRepository(conf *config.Store, storeClient storeClient.Client, httpClient httpClient.Client) Repository {
	return &repositoryImpl{
		conf:        conf,
		storeClient: storeClient,
		httpClient:  httpClient,
	}
}

func (r *repositoryImpl) Upload(file []byte, bucketName string, objectKey string) (url string, key string, err error) {
	ctx := context.Background()
	_, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	buffer := bytes.NewReader(file)

	uploadOutput, err := r.storeClient.PutObject(ctx, bucketName, objectKey, buffer,
		buffer.Size(), minio.PutObjectOptions{})
	if err != nil {
		return "", "", errors.Wrap(err, fmt.Sprintf("Couldn't upload object to %v/%v.", bucketName, objectKey))
	}

	return r.GetURL(bucketName, objectKey), uploadOutput.Key, nil
}

func (r *repositoryImpl) Delete(bucketName string, objectKey string) (err error) {
	ctx := context.Background()
	_, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err = r.storeClient.RemoveObject(ctx, bucketName, objectKey, opts)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Couldn't delete object %v/%v.", bucketName, objectKey))
	}

	return nil
}

func (r *repositoryImpl) Get(bucketName string, objectKey string) (url string, err error) {
	ctx := context.Background()
	_, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	url = r.GetURL(bucketName, objectKey)

	resp, err := r.httpClient.Get(url)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Couldn't get object %v/%v.", bucketName, objectKey))
	}
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	return url, nil
}

func (r *repositoryImpl) GetURL(bucketName string, objectKey string) string {
	return "https://" + r.conf.Endpoint + "/" + bucketName + "/" + objectKey
}
