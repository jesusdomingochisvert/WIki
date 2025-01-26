package repositories

import (
	"context"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/entities"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/interfaces"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type userRepository struct {
	collection interfaces.Collection
}

func NewUserRepository(db interfaces.Database) repositories.UserRepository {
	return &userRepository{collection: db.Collection("users")}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []entities.User
	for cur.Next(ctx) {
		var user entities.User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return users, nil
}
