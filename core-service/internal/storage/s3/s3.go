package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3PhotoService - реализация PhotoService, которая загружает фото на S3
type S3Storage struct {
	s3Client *s3.S3
	bucket   string
}

// NewS3Storage создает новый сервис для загрузки фотографий на S3
func NewS3Storage(bucket string, s3Client *s3.S3) *S3Storage {
	return &S3Storage{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

// UploadPhoto загружает фото на S3
func (s *S3Storage) UploadPhoto(ctx context.Context, file multipart.File, fileSize int64, fileName string) (string, error) {

	// Чтение файла для передачи на S3
	fileBytes := make([]byte, fileSize)
	if _, err := file.Read(fileBytes); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Подготовка параметров для загрузки на S3
	uploadParams := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filepath.Base(fileName)),
		Body:   file,
	}

	// Загрузка файла на S3
	_, err := s.s3Client.PutObjectWithContext(ctx, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Возвращаем URL загруженного файла
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, fileName)
	return fileURL, nil
}
