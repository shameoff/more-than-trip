package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
	"github.com/shameoff/more-than-trip/core/internal/lib/logger/sl"
)

type CoreService interface {
	UploadPhoto(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, metadata models.Photo) error
	// Работа с фотографиями
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
	GetRegions(ctx context.Context) ([]models.Region, error)
	UpdateRegion(ctx context.Context, regionId uuid.UUID, data models.Region) error
}

type CoreHandler struct {
	service CoreService
	logger  *slog.Logger
}

func NewCoreHandler(service CoreService, logger *slog.Logger) *CoreHandler {
	return &CoreHandler{
		service: service,
		logger:  logger,
	}
}

// UploadPhoto - HTTP handler для загрузки фото и метаданных
func (h *CoreHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	// Получаем файл из запроса
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get file from request", sl.Err(err))
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем метаданные (JSON) из текстового поля
	metadataStr := r.FormValue("metadata")
	if metadataStr == "" {
		h.logger.Error("missing metadata")
		http.Error(w, "missing metadata", http.StatusBadRequest)
		return
	}

	// Парсим метаданные в структуру
	metadata := models.Photo{}
	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		h.logger.Error("failed to parse metadata", sl.Err(err))
		http.Error(w, "invalid metadata format", http.StatusBadRequest)
		return
	}

	// Передача файла и метаданных на уровень бизнес-логики
	err = h.service.UploadPhoto(ctx, file, fileHeader, metadata)
	if err != nil {
		h.logger.Error("failed to upload file via service", sl.Err(err))
		http.Error(w, "file upload failed", http.StatusInternalServerError)
		return
	}

	h.logger.Info("file uploaded successfully")
	w.WriteHeader(http.StatusOK)
}

func (h *CoreHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	// Получаем photoId из URL параметров
	photoIdStr := chi.URLParam(r, "photoId")
	photoId, err := uuid.Parse(photoIdStr)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	// Вызываем бизнес-логику для удаления фотографии
	err = h.service.DeletePhoto(r.Context(), photoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Photo deleted successfully"))
}
func (h *CoreHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	photoIdStr := chi.URLParam(r, "photoId")
	photoId, err := uuid.Parse(photoIdStr)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	photo, err := h.service.GetPhoto(r.Context(), photoId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photo)
}
func (h *CoreHandler) GetPhotos(w http.ResponseWriter, r *http.Request) {
	var filters models.PhotoFiltersDTO

	// Пример получения фильтров из запроса
	if regionId := r.URL.Query().Get("region"); regionId != "" {
		filters.RegionId, _ = uuid.Parse(regionId)
	}
	if tag := r.URL.Query().Get("tag"); tag != "" {
		filters.TagKey = tag
	}

	photos, err := h.service.GetPhotos(r.Context(), filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get photos: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}
func (h *CoreHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	photoIdStr := chi.URLParam(r, "photoId")
	photoId, err := uuid.Parse(photoIdStr)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	var updatedPhoto models.Photo
	if err := json.NewDecoder(r.Body).Decode(&updatedPhoto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdatePhoto(r.Context(), photoId, updatedPhoto)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Photo updated successfully"))
}

func (h *CoreHandler) LikePhoto(w http.ResponseWriter, r *http.Request) {
	photoIdStr := chi.URLParam(r, "photoId")
	photoId, err := uuid.Parse(photoIdStr)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	userIdStr := r.Header.Get("User-Id") // Предположим, что ID пользователя передается в заголовке
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.service.LikePhoto(r.Context(), photoId, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to like photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Photo liked successfully"))
}
func (h *CoreHandler) DislikePhoto(w http.ResponseWriter, r *http.Request) {
	photoIdStr := chi.URLParam(r, "photoId")
	photoId, err := uuid.Parse(photoIdStr)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	userIdStr := r.Header.Get("User-Id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.service.DislikePhoto(r.Context(), photoId, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to dislike photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Photo disliked successfully"))
}

func (h *CoreHandler) CreateRegion(w http.ResponseWriter, r *http.Request) {
	var region models.Region
	if err := json.NewDecoder(r.Body).Decode(&region); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.CreateRegion(r.Context(), region)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create region: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Region created successfully"))
}
func (h *CoreHandler) DeleteRegion(w http.ResponseWriter, r *http.Request) {
	regionIdStr := chi.URLParam(r, "regionId")
	regionId, err := uuid.Parse(regionIdStr)
	if err != nil {
		http.Error(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteRegion(r.Context(), regionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete region: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Region deleted successfully"))
}
func (h *CoreHandler) GetRegion(w http.ResponseWriter, r *http.Request) {
	regionIdStr := chi.URLParam(r, "regionId")
	regionId, err := uuid.Parse(regionIdStr)
	if err != nil {
		http.Error(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	region, err := h.service.GetRegionById(r.Context(), regionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get region: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(region)
}
func (h *CoreHandler) GetRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := h.service.GetRegions(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get regions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regions)
}
func (h *CoreHandler) UpdateRegion(w http.ResponseWriter, r *http.Request) {
	regionIdStr := chi.URLParam(r, "regionId")
	regionId, err := uuid.Parse(regionIdStr)
	if err != nil {
		http.Error(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	var updatedRegion models.Region
	if err := json.NewDecoder(r.Body).Decode(&updatedRegion); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateRegion(r.Context(), regionId, updatedRegion)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update region: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Region updated successfully"))
}

func (h *CoreHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.CreateTag(r.Context(), tag)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create tag: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Tag created successfully"))
}
func (h *CoreHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.GetTags(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get tags: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}
func (h *CoreHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagId := chi.URLParam(r, "tagId")

	err := h.service.DeleteTag(r.Context(), tagId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete tag: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tag deleted successfully"))
}

// ============================NOT IMPLEMENTED=================================
// Заглушка для всех хендлеров, возвращающих 501
func notImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("This functionality is not yet implemented"))
}

// Реализация хендлеров для работы с вызовами
func (h *CoreHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) UpdateChallenge(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) GetChallenge(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) GetChallenges(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) GetMyChallenges(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

// Реализация хендлеров для работы с пользователями
func (h *CoreHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

// Реализация хендлеров для работы с поездками
func (h *CoreHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) GetTrip(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}

func (h *CoreHandler) GetTrips(w http.ResponseWriter, r *http.Request) {
	notImplemented(w)
}
