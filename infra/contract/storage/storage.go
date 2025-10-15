package storage

import (
	"context"
	"time"
)

type Storage interface {
	// PutObject puts the object with the specified key.
	PutObject(ctx context.Context, objectKey string, content []byte, opts ...PutOptFn) error
	// GetObject returns the object with the specified key.
	GetObject(ctx context.Context, objectKey string) ([]byte, error)
	// DeleteObject deletes the object with the specified key.
	DeleteObject(ctx context.Context, objectKey string) error
	// GetObjectUrl returns a presigned URL for the object.
	// The URL is valid for the specified duration.
	GetObjectUrl(ctx context.Context, objectKey string, opts ...GetOptFn) (string, error)
}

type FileInfo struct {
	Key          string            `json:"key"`
	LastModified time.Time         `json:"last_modified"`
	ETag         string            `json:"etag"`
	Size         int64             `json:"size"`
	URL          string            `json:"url"`
	Tagging      map[string]string `json:"tagging"`
}
