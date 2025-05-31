package storage

import (
	"context"
	"io"
	"time"
)

// ObjectStorage defines interface for cloud storage operations
type ObjectStorage interface {
	// PutObject uploads an object to storage
	PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) error

	// GetObject downloads an object from storage
	GetObject(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error)

	// DeleteObject removes an object from storage
	DeleteObject(ctx context.Context, bucketName string, objectName string) error

	// ListObjects lists objects in a bucket with optional prefix
	ListObjects(ctx context.Context, bucketName string, prefix string) ([]ObjectInfo, error)

	// BucketExists checks if a bucket exists
	BucketExists(ctx context.Context, bucketName string) (bool, error)

	// MakeBucket creates a new bucket
	MakeBucket(ctx context.Context, bucketName string) error

	// RemoveBucket removes a bucket
	RemoveBucket(ctx context.Context, bucketName string) error

	// GetObjectURL returns a URL for accessing the object
	GetObjectURL(ctx context.Context, bucketName string, objectName string, expires time.Duration) (string, error)

	// GetObjectURI returns a URI for accessing the object
	GetObjectURI(ctx context.Context, bucketName string, objectName string) (string, error)
}

// ObjectInfo contains information about an object in storage
type ObjectInfo struct {
	Key          string    // Object name
	Size         int64     // Object size
	LastModified time.Time // Last modified time
	ETag         string    // ETag of the object
	ContentType  string    // Content type of the object
}
