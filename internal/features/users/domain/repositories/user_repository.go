package repositories

import (
	"context"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/entities"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]entities.User, error)
}
