package core

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes регистрирует все маршруты для вашего HTTP API.
// Вызывается после инициализации хендлеров(контроллеров)
func RegisterRoutes(router *chi.Mux, coreHandler *CoreHandler) {
	// Эндпоинты для загрузки и удаления фото
	router.Post("/api/photo", coreHandler.UploadPhoto)
	router.Get("/api/photo/{UUID}", coreHandler.GetPhoto)
	router.Get("/api/photos", coreHandler.GetPhotos)
	router.Delete("/api/photo/{UUID}", coreHandler.DeletePhoto)
	router.Put("/api/photo/{UUID}", coreHandler.UpdatePhoto)
	router.Post("/api/photo/{UUID}/like", coreHandler.LikePhoto)
	router.Post("/api/photo/{UUID}/dislike", coreHandler.DislikePhoto)

	// Маршруты для работы с регионами
	router.Post("/api/region", coreHandler.CreateRegion)
	router.Delete("/api/region/{UUID}", coreHandler.DeleteRegion)
	router.Get("/api/region/{UUID}", coreHandler.GetRegion)
	router.Get("/api/regions", coreHandler.GetRegions)
	router.Put("/api/region/{UUID}", coreHandler.UpdateRegion)

	// Маршруты для работы с тегами
	router.Post("/api/tag", coreHandler.CreateTag)
	router.Get("/api/tags", coreHandler.GetTags)
	router.Delete("/api/tag/{UUID}", coreHandler.DeleteTag)

	// Маршруты для работы с вызовами (challenges)
	router.Post("/api/challenge", coreHandler.CreateChallenge)
	router.Delete("/api/challenge/{UUID}", coreHandler.DeleteChallenge)
	router.Put("/api/challenge/{UUID}", coreHandler.UpdateChallenge)
	router.Get("/api/challenge/{UUID}", coreHandler.GetChallenge)
	router.Get("/api/challenges", coreHandler.GetChallenges)
	router.Get("/api/mychallenges", coreHandler.GetMyChallenges)

	// Маршруты для работы с пользователями
	router.Post("/api/user", coreHandler.CreateUser)
	router.Delete("/api/user/{UUID}", coreHandler.DeleteUser)
	router.Put("/api/user/{UUID}", coreHandler.UpdateUser)
	router.Get("/api/user/{UUID}", coreHandler.GetUser)
	router.Get("/api/users", coreHandler.GetUsers)

	// Маршруты для работы с поездками
	router.Post("/api/trip", coreHandler.CreateTrip)
	router.Delete("/api/trip/{UUID}", coreHandler.DeleteTrip)
	router.Put("/api/trip/{UUID}", coreHandler.UpdateTrip)
	router.Get("/api/trip/{UUID}", coreHandler.GetTrip)
	router.Get("/api/trips", coreHandler.GetTrips)

}
