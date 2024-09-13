// Бизнес логика, связанная с пользователем
package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
)

// Работа с пользователями
func (c *CoreService) CreateUser(ctx context.Context, user models.User) error {
	return c.storage.CreateUser(ctx, user)
}

func (c *CoreService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	return c.storage.DeleteUser(ctx, userId)
}

func (c *CoreService) GetUserById(ctx context.Context, userId uuid.UUID) (models.User, error) {
	return c.storage.GetUserById(ctx, userId)
}

func (c *CoreService) GetUsers(ctx context.Context) ([]models.User, error) {
	return c.storage.GetUsers(ctx)
}

func (c *CoreService) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return c.storage.GetUserByUsername(ctx, username)
}

func (c *CoreService) UpdateUser(ctx context.Context, userId uuid.UUID, data models.User) error {
	return c.storage.UpdateUser(ctx, userId, data)
}
