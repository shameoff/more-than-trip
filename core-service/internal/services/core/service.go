// Здесь должна быть бизнес логика и зависимость от Storage
package core

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
)

type CoreStorage interface {
	SavePhoto(ctx context.Context, photoData models.Photo) (string, error)
	GetPhoto(ctx context.Context, photoID uuid.UUID) (models.Photo, error)
	GetPhotosByTripID(ctx context.Context, tripID uuid.UUID) ([]models.Photo, error)
}

// S3PhotoService - интерфейс для работы с файлами
type S3PhotoStorage interface {
	UploadPhoto(ctx context.Context, file multipart.File, fileSize int64, fileName string) (string, error)
}

type CoreService struct {
	uploadDir string
	log       *slog.Logger
	storage   CoreStorage
	s3Storage S3PhotoStorage
}

// NewCoreService создает сервис для загрузки файлов в локальную файловую систему
func NewCoreService(uploadDir string,
	log *slog.Logger,
	storage CoreStorage,
	s3storage S3PhotoStorage,
) *CoreService {
	return &CoreService{
		uploadDir: uploadDir,
		log:       log,
		storage:   storage,
		s3Storage: s3storage,
	}
}

func (s *CoreService) UploadPhoto(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, metadata models.Photo) (string, error) {
	// Генерация уникального имени файла
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)

	// Загрузка файла на S3
	fileURL, err := s.s3Storage.UploadPhoto(ctx, file, fileHeader.Size, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	metadata.ImgUrl = fileURL
	_, err = s.storage.SavePhoto(ctx, metadata)
	if err != nil {
		return "", fmt.Errorf("failed to save photo metadata: %w", err)
	}

	return fileURL, nil
}

func (s *CoreService) DownloadPhoto(ctx context.Context, fileID string) error {
	filePath := filepath.Join(s.uploadDir, fileID)
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return os.ErrNotExist
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
