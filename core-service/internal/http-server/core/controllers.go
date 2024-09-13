package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
	"github.com/shameoff/more-than-trip/core/internal/lib/logger/sl"
)

type CoreService interface {
	UploadPhoto(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, metadata models.Photo) (string, error)
	DownloadPhoto(ctx context.Context, fileID string) error
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
	photoURL, err := h.service.UploadPhoto(ctx, file, fileHeader, metadata)
	if err != nil {
		h.logger.Error("failed to upload file via service", sl.Err(err))
		http.Error(w, "file upload failed", http.StatusInternalServerError)
		return
	}

	h.logger.Info("file uploaded successfully", slog.String("file_url", photoURL))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", photoURL)))
}

// DownloadPhoto - HTTP handler для скачивания фото
func (h *CoreHandler) DownloadPhoto(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	fileID := chi.URLParam(r, "fileID")
	if fileID == "" {
		h.logger.Error("missing fileID parameter")
		http.Error(w, "fileID is required", http.StatusBadRequest)
		return
	}

	err := h.service.DownloadPhoto(ctx, fileID)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.logger.Error("file not found", slog.String("fileID", fileID))
			http.Error(w, "file not found", http.StatusNotFound)
		} else {
			h.logger.Error("failed to download file", sl.Err(err))
			http.Error(w, "failed to download file", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("file downloaded successfully", slog.String("fileID", fileID))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File downloaded successfully"))
}

func (h *CoreHandler) DeletePhoto(w http.ResponseWriter, r *http.Request)  {}
func (h *CoreHandler) GetPhoto(w http.ResponseWriter, r *http.Request)     {}
func (h *CoreHandler) GetPhotos(w http.ResponseWriter, r *http.Request)    {}
func (h *CoreHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request)  {}
func (h *CoreHandler) TagPhoto(w http.ResponseWriter, r *http.Request)     {}
func (h *CoreHandler) LikePhoto(w http.ResponseWriter, r *http.Request)    {}
func (h *CoreHandler) DislikePhoto(w http.ResponseWriter, r *http.Request) {}

func (h *CoreHandler) CreateRegion(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) DeleteRegion(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) GetRegion(w http.ResponseWriter, r *http.Request)    {}
func (h *CoreHandler) GetRegions(w http.ResponseWriter, r *http.Request)   {}
func (h *CoreHandler) UpdateRegion(w http.ResponseWriter, r *http.Request) {}

func (h *CoreHandler) CreateTag(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) GetTags(w http.ResponseWriter, r *http.Request)   {}
func (h *CoreHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {}

func (h *CoreHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) DeleteChallenge(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) UpdateChallenge(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) GetChallenge(w http.ResponseWriter, r *http.Request)    {}
func (h *CoreHandler) GetChallenges(w http.ResponseWriter, r *http.Request)   {}
func (h *CoreHandler) GetMyChallenges(w http.ResponseWriter, r *http.Request) {}

func (h *CoreHandler) CreateUser(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) GetUser(w http.ResponseWriter, r *http.Request)    {}
func (h *CoreHandler) GetUsers(w http.ResponseWriter, r *http.Request)   {}

func (h *CoreHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) DeleteTrip(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) UpdateTrip(w http.ResponseWriter, r *http.Request) {}
func (h *CoreHandler) GetTrip(w http.ResponseWriter, r *http.Request)    {}
func (h *CoreHandler) GetTrips(w http.ResponseWriter, r *http.Request)   {}
