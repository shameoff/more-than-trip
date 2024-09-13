// Работа с бизнес-логикой поездок
package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
)

// Работа с поездками
func (c *CoreService) CreateTrip(ctx context.Context, trip models.Trip) error {
	return c.storage.CreateTrip(ctx, trip)
}

func (c *CoreService) DeleteTrip(ctx context.Context, tripId uuid.UUID) error {
	return c.storage.DeleteTrip(ctx, tripId)
}

func (c *CoreService) GetTripById(ctx context.Context, tripId uuid.UUID) (models.Trip, error) {
	return c.storage.GetTripById(ctx, tripId)
}

func (c *CoreService) GetTripsByUserId(ctx context.Context, userId uuid.UUID) ([]models.Trip, error) {
	return c.storage.GetTripsByUserId(ctx, userId)
}

func (c *CoreService) GetTripsByRegionId(ctx context.Context, regionId uuid.UUID) ([]models.Trip, error) {
	return c.storage.GetTripsByRegionId(ctx, regionId)
}

func (c *CoreService) GetTripsByTag(ctx context.Context, tagId string) ([]models.Trip, error) {
	return c.storage.GetTripsByTag(ctx, tagId)
}

func (c *CoreService) UpdateTrip(ctx context.Context, tripId uuid.UUID, data models.Trip) error {
	return c.storage.UpdateTrip(ctx, tripId, data)
}
