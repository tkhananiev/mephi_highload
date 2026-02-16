package services

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IntegrationService struct {
	client *minio.Client
	bucket string
}

func NewIntegrationService(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*IntegrationService, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &IntegrationService{client: cli, bucket: bucket}, nil
}

func (s *IntegrationService) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
}

func (s *IntegrationService) UploadText(ctx context.Context, objectName string, content []byte) (string, error) {
	r := bytes.NewReader(content)
	_, err := s.client.PutObject(ctx, s.bucket, objectName, r, int64(len(content)), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("s3://%s/%s", s.bucket, objectName), nil
}

func DefaultAuditObjectName() string {
	return fmt.Sprintf("audit-%s.log", time.Now().UTC().Format("20060102T150405Z"))
}
