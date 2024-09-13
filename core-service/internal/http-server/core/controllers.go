package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shameoff/more-than-trip/core/internal/lib/logger/sl"
)

type CoreService interface {
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	DeleteFile(ctx context.Context, fileID string) error
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

// UploadFile - HTTP handler для загрузки файлов
func (h *CoreHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get file from request", sl.Err(err))
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadedFilePath, err := h.service.UploadFile(ctx, file, fileHeader)
	if err != nil {
		h.logger.Error("failed to upload file", sl.Err(err))
		http.Error(w, "file upload failed", http.StatusInternalServerError)
		return
	}

	h.logger.Info("file uploaded successfully", slog.String("file_path", uploadedFilePath))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", uploadedFilePath)))
}

// DeleteFile - HTTP handler для удаления файлов
func (h *CoreHandler) DownloadPhoto(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	fileID := chi.URLParam(r, "fileID")
	if fileID == "" {
		h.logger.Error("missing fileID parameter")
		http.Error(w, "fileID is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteFile(ctx, fileID)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.logger.Error("file not found", slog.String("fileID", fileID))
			http.Error(w, "file not found", http.StatusNotFound)
		} else {
			h.logger.Error("failed to delete file", sl.Err(err))
			http.Error(w, "failed to delete file", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("file deleted successfully", slog.String("fileID", fileID))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
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
