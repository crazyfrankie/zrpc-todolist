package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/storage"
)

type minioStore struct {
	bucketName string
	client     *minio.Client
}

func getMinIOClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func New(ctx context.Context, endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (storage.Storage, error) {
	m, err := getMinIOClient(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return &minioStore{client: m, bucketName: bucketName}, nil
}

func (m *minioStore) PutObject(ctx context.Context, objectKey string, content []byte, opts ...storage.PutOptFn) error {
	option := storage.PutOption{}
	for _, opt := range opts {
		opt(&option)
	}

	minioOpts := minio.PutObjectOptions{}
	if option.ContentType != nil {
		minioOpts.ContentType = *option.ContentType
	}

	if option.ContentEncoding != nil {
		minioOpts.ContentEncoding = *option.ContentEncoding
	}

	if option.ContentDisposition != nil {
		minioOpts.ContentDisposition = *option.ContentDisposition
	}

	if option.ContentLanguage != nil {
		minioOpts.ContentLanguage = *option.ContentLanguage
	}

	if option.Expires != nil {
		minioOpts.Expires = *option.Expires
	}

	_, err := m.client.PutObject(ctx, m.bucketName, objectKey,
		bytes.NewReader(content), int64(len(content)), minioOpts)
	if err != nil {
		return fmt.Errorf("PutObject failed: %v", err)
	}
	return nil
}

func (m *minioStore) GetObject(ctx context.Context, objectKey string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}
	defer obj.Close()
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("ReadObject failed: %v", err)
	}
	return data, nil
}

func (m *minioStore) DeleteObject(ctx context.Context, objectKey string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("DeleteObject failed: %v", err)
	}
	return nil
}

func (m *minioStore) GetObjectUrl(ctx context.Context, objectKey string, opts ...storage.GetOptFn) (string, error) {
	option := storage.GetOption{}
	for _, opt := range opts {
		opt(&option)
	}

	if option.Expire == 0 {
		option.Expire = 3600 * 24 * 7
	}

	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectKey, time.Duration(option.Expire)*time.Second, reqParams)
	if err != nil {
		return "", fmt.Errorf("GetObjectUrl failed: %v", err)
	}

	return presignedURL.String(), nil
}
