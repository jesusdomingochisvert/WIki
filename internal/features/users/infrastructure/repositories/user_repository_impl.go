package repositories

import (
	"context"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/entities"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/interfaces"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/repositories"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/infrastructure/db/schema"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/infrastructure/mappers"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type userRepository struct {
	collection interfaces.Collection
}

func NewUserRepository(db interfaces.Database) repositories.UserRepository {
	return &userRepository{collection: db.Collection("users")}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]entities.UserEntity, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var userEntities []entities.UserEntity
	for cur.Next(ctx) {
		var userSchema schema.UserSchema
		if err := cur.Decode(&userSchema); err != nil {
			return nil, err
		}
		userEntity := mappers.FromUserSchema(userSchema)
		userEntities = append(userEntities, userEntity)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return userEntities, nil
}
