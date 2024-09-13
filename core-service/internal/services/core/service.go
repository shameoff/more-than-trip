// Здесь должна быть бизнес логика и зависимость от Storage
package core

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type CoreStorage interface {
}
type CoreService struct {
	uploadDir string
	log       *slog.Logger
	storage   CoreStorage
}

// NewLocalFileUploadService создает сервис для загрузки файлов в локальную файловую систему
func NewCoreService(uploadDir string,
	log *slog.Logger,
	storage CoreStorage,
) *CoreService {
	return &CoreService{
		uploadDir: uploadDir,
		log:       log,
		storage:   storage,
	}
}

func (s *CoreService) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	fileID := uuid.New().String()
	filePath := filepath.Join(s.uploadDir, fileID+"_"+fileHeader.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	return filePath, nil
}

func (s *CoreService) DeleteFile(ctx context.Context, fileID string) error {
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
