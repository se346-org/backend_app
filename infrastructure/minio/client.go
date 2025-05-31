package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/chat-socio/backend/configuration"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/chat-socio/backend/pkg/storage"
)

type minioClient struct {
	client *minio.Client
	obs    *observability.Observability
}

func NewMinioClient(cfg *configuration.MinioConfig, obs *observability.Observability) (storage.ObjectStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, cfg.Token),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &minioClient{client: client, obs: obs}, nil
}

func (m *minioClient) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) error {
	_, err := m.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

func (m *minioClient) GetObject(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *minioClient) DeleteObject(ctx context.Context, bucketName string, objectName string) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m *minioClient) ListObjects(ctx context.Context, bucketName string, prefix string) ([]storage.ObjectInfo, error) {
	objectCh := m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var objects []storage.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, storage.ObjectInfo{
			Key:          object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ETag:         object.ETag,
			ContentType:  object.ContentType,
		})
	}
	return objects, nil
}

func (m *minioClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return m.client.BucketExists(ctx, bucketName)
}

func (m *minioClient) MakeBucket(ctx context.Context, bucketName string) error {
	return m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

func (m *minioClient) RemoveBucket(ctx context.Context, bucketName string) error {
	return m.client.RemoveBucket(ctx, bucketName)
}

func (m *minioClient) GetObjectURL(ctx context.Context, bucketName string, objectName string, expires time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *minioClient) GetObjectURI(ctx context.Context, bucketName string, objectName string) (string, error) {
	return fmt.Sprintf("/%s/%s", bucketName, objectName), nil
}
