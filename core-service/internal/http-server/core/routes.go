package core

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes регистрирует все маршруты для вашего HTTP API.
// Вызывается после инициализации хендлеров(контроллеров)
func RegisterRoutes(router *chi.Mux, coreHandler *CoreHandler) {
	// Эндпоинты для загрузки и удаления файлов
	router.Post("/upload", coreHandler.UploadPhoto)
	router.Delete("/files/{fileID}", coreHandler.DeletePhoto)

	// Вы можете добавлять и другие группы маршрутов здесь
	// Например, эндпоинты для профилей, авторизации и т.д.
}
