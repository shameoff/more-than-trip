// Бизнес-логика связанная с объектами фотографий и общие элементы пакета
package core

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
)

type CoreStorage interface {
	// Работа с фотографиями
	SavePhoto(ctx context.Context, data models.Photo) error
	GetPhoto(ctx context.Context, photoId uuid.UUID) (models.Photo, error)
	GetPhotos(ctx context.Context, filters models.PhotoFiltersDTO) ([]models.Photo, error)
	UpdatePhoto(ctx context.Context, photoId uuid.UUID, data models.Photo) error
	DeletePhoto(ctx context.Context, photoId uuid.UUID) error
	GetPhotosByTripId(ctx context.Context, tripId uuid.UUID) ([]models.Photo, error)
	GetPhotosByUserId(ctx context.Context, userId uuid.UUID) ([]models.Photo, error)
	GetPhotosByRegionId(ctx context.Context, regionId uuid.UUID) ([]models.Photo, error)

	// Работа с лайками
	LikePhoto(ctx context.Context, photoId uuid.UUID, userId uuid.UUID) error
	DislikePhoto(ctx context.Context, photoId uuid.UUID, userId uuid.UUID) error

	// Работа с тегами
	CreateTag(ctx context.Context, tag models.Tag) error
	DeleteTag(ctx context.Context, tagId string) error
	GetTags(ctx context.Context) ([]models.Tag, error)
	GetTagsByPhotoId(ctx context.Context, photoId uuid.UUID) ([]models.Tag, error)

	// Работа с регионами
	CreateRegion(ctx context.Context, region models.Region) error
	DeleteRegion(ctx context.Context, regionId uuid.UUID) error
	GetRegionById(ctx context.Context, regionId uuid.UUID) (models.Region, error)
	GetRegionByKey(ctx context.Context, regionKey string) (models.Region, error)
	GetRegions(ctx context.Context) ([]models.Region, error)
	UpdateRegion(ctx context.Context, regionId uuid.UUID, data models.Region) error

	// Работа с поездками
	CreateTrip(ctx context.Context, trip models.Trip) error
	DeleteTrip(ctx context.Context, tripId uuid.UUID) error
	GetTripById(ctx context.Context, tripId uuid.UUID) (models.Trip, error)
	GetTripsByUserId(ctx context.Context, userId uuid.UUID) ([]models.Trip, error)
	GetTripsByRegionId(ctx context.Context, regionId uuid.UUID) ([]models.Trip, error)
	GetTripsByTag(ctx context.Context, tagId string) ([]models.Trip, error)
	UpdateTrip(ctx context.Context, tripId uuid.UUID, data models.Trip) error

	// Работа с пользователями
	CreateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, userId uuid.UUID) error
	GetUserById(ctx context.Context, userId uuid.UUID) (models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, data models.User) error
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

func (s *CoreService) UploadPhoto(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, metadata models.Photo) error {
	// Генерация уникального имени файла
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)

	// Загрузка файла на S3
	fileURL, err := s.s3Storage.UploadPhoto(ctx, file, fileHeader.Size, fileName)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	metadata.ImgUrl = fileURL
	// Загрузка в БД
	err = s.storage.SavePhoto(ctx, metadata)
	if err != nil {
		return fmt.Errorf("failed to save photo metadata: %w", err)
	}

	return nil
}

func (c *CoreService) GetPhoto(ctx context.Context, photoId uuid.UUID) (models.Photo, error) {
	photo, err := c.storage.GetPhoto(ctx, photoId)
	if err != nil {
		// Логирование ошибки, если нужно
		return models.Photo{}, fmt.Errorf("failed to get photo: %w", err)
	}
	return photo, nil
}

func (c *CoreService) GetPhotos(ctx context.Context, filters models.PhotoFiltersDTO) ([]models.Photo, error) {
	photos, err := c.storage.GetPhotos(ctx, filters)
	if err != nil {
		// Логирование ошибки, если нужно
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}
	return photos, nil
}

func (c *CoreService) UpdatePhoto(ctx context.Context, photoId uuid.UUID, data models.Photo) error {
	err := c.storage.UpdatePhoto(ctx, photoId, data)
	if err != nil {
		// Логирование ошибки, если нужно
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}

func (c *CoreService) DeletePhoto(ctx context.Context, photoId uuid.UUID) error {
	err := c.storage.DeletePhoto(ctx, photoId)
	if err != nil {
		// Логирование ошибки, если нужно
		return fmt.Errorf("failed to delete photo: %w", err)
	}
	return nil
}

func (c *CoreService) GetPhotosByTripId(ctx context.Context, tripId uuid.UUID) ([]models.Photo, error) {
	photos, err := c.storage.GetPhotosByTripId(ctx, tripId)
	if err != nil {
		// Логирование ошибки, если нужно
		return nil, fmt.Errorf("failed to get photos by trip ID: %w", err)
	}
	return photos, nil
}

func (c *CoreService) GetPhotosByUserId(ctx context.Context, userId uuid.UUID) ([]models.Photo, error) {
	photos, err := c.storage.GetPhotosByUserId(ctx, userId)
	if err != nil {
		// Логирование ошибки, если нужно
		return nil, fmt.Errorf("failed to get photos by user ID: %w", err)
	}
	return photos, nil
}

func (c *CoreService) GetPhotosByRegionId(ctx context.Context, regionId uuid.UUID) ([]models.Photo, error) {
	photos, err := c.storage.GetPhotosByRegionId(ctx, regionId)
	if err != nil {
		// Логирование ошибки, если нужно
		return nil, fmt.Errorf("failed to get photos by region ID: %w", err)
	}
	return photos, nil
}

func (c *CoreService) LikePhoto(ctx context.Context, photoId uuid.UUID, userId uuid.UUID) error {
	err := c.storage.LikePhoto(ctx, photoId, userId)
	if err != nil {
		// Логирование ошибки, если нужно
		return fmt.Errorf("failed to like photo: %w", err)
	}
	return nil
}

func (c *CoreService) DislikePhoto(ctx context.Context, photoId uuid.UUID, userId uuid.UUID) error {
	err := c.storage.DislikePhoto(ctx, photoId, userId)
	if err != nil {
		// Логирование ошибки, если нужно
		return fmt.Errorf("failed to dislike photo: %w", err)
	}
	return nil
}

// Работа с тегами
func (c *CoreService) CreateTag(ctx context.Context, tag models.Tag) error {
	err := c.storage.CreateTag(ctx, tag)
	if err != nil {
		// Логирование ошибки, если нужно
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (c *CoreService) DeleteTag(ctx context.Context, tagId string) error {
	err := c.storage.DeleteTag(ctx, tagId)
	if err != nil {
		// Логирование ошибки, если нужно
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

func (c *CoreService) GetTags(ctx context.Context) ([]models.Tag, error) {
	tags, err := c.storage.GetTags(ctx)
	if err != nil {
		// Логирование ошибки, если нужно
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return tags, nil
}

func (c *CoreService) GetTagsByPhotoId(ctx context.Context, photoId uuid.UUID) ([]models.Tag, error) {
	tags, err := c.storage.GetTagsByPhotoId(ctx, photoId)
	if err != nil {
		// Логирование ошибки, если нужно
		return nil, fmt.Errorf("failed to get tags by photo ID: %w", err)
	}
	return tags, nil
}
