// Бизнес-логика связанная с регионами
package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
)

func (c *CoreService) CreateRegion(ctx context.Context, region models.Region) error {
	return c.storage.CreateRegion(ctx, region)
}

func (c *CoreService) DeleteRegion(ctx context.Context, regionId uuid.UUID) error {
	return c.storage.DeleteRegion(ctx, regionId)
}

func (c *CoreService) GetRegionById(ctx context.Context, regionId uuid.UUID) (models.Region, error) {
	return c.storage.GetRegionById(ctx, regionId)
}

func (c *CoreService) GetRegionByKey(ctx context.Context, regionKey string) (models.Region, error) {
	return c.storage.GetRegionByKey(ctx, regionKey)
}

func (c *CoreService) GetRegions(ctx context.Context) ([]models.Region, error) {
	return c.storage.GetRegions(ctx)
}

func (c *CoreService) UpdateRegion(ctx context.Context, regionId uuid.UUID, data models.Region) error {
	return c.storage.UpdateRegion(ctx, regionId, data)
}
