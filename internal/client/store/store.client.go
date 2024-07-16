package store

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	PutObject(ctx context.Context, bucketName string, objectName string, reader *bytes.Reader) (info *s3.PutObjectOutput, err error)
	RemoveObject(ctx context.Context, bucketName string, objectName string) error
}

type clientImpl struct {
	*s3.S3
}

func NewClient(s3Client *s3.S3) Client {
	return &clientImpl{s3Client}
}

func (c *clientImpl) PutObject(ctx context.Context, bucketName string, objectName string, reader *bytes.Reader) (info *s3.PutObjectOutput, err error) {
	return c.S3.PutObject(&s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectName,
		Body:   reader,
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
	})
}

func (c *clientImpl) RemoveObject(ctx context.Context, bucketName string, objectName string) error {
	_, err := c.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectName,
	})
	return err
}
