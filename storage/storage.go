package storage

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

var (
	ErrUploadFailed    = errors.New("failed to upload file")
	ErrDeleteFailed    = errors.New("failed to delete file")
	ErrDownloadFailed  = errors.New("failed to download file")
	ErrInvalidFilePath = errors.New("invalid file path")
)

type Storage interface {
	Upload(relativePath string, file io.Reader, contentType string) (string, error)
	Download(publicURL string) ([]byte, string, error)
	Delete(publicURL string) error
}

type StorageClient struct {
	client *storage_go.Client
	bucket string
}

func NewStorageClient(baseURL string, serviceKey string, bucket string) Storage {
	return &StorageClient{
		client: storage_go.NewClient(baseURL, serviceKey, map[string]string{
			"apiKey": serviceKey,
		}),
		bucket: bucket,
	}
}

func (s *StorageClient) Upload(relativePath string, file io.Reader, contentType string) (string, error) {
	upsert := true
	options := storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &upsert,
	}

	_, err := s.client.UploadFile(s.bucket, relativePath, file, options)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	return s.publicURL(relativePath), nil
}

func (s *StorageClient) Delete(publicURL string) error {
	relativePath := s.pathFromURL(publicURL)
	if relativePath == "" {
		return nil
	}

	_, err := s.client.RemoveFile(s.bucket, []string{relativePath})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	return nil
}

func (s *StorageClient) Download(publicURL string) ([]byte, string, error) {
	relativePath := s.pathFromURL(publicURL)
	if relativePath == "" {
		return nil, "", ErrInvalidFilePath
	}

	data, err := s.client.DownloadFile(s.bucket, relativePath)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}

	return data, ContentType(relativePath), nil
}

func ContentType(relativePath string) string {
	switch strings.ToLower(filepath.Ext(relativePath)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func (s *StorageClient) publicURL(relativePath string) string {
	return s.client.GetPublicUrl(s.bucket, relativePath).SignedURL
}

func (s *StorageClient) pathFromURL(publicURL string) string {
	prefix := "/object/public/" + s.bucket + "/"
	index := strings.Index(publicURL, prefix)
	if index < 0 {
		return ""
	}
	return publicURL[index+len(prefix):]
}
