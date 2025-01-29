package mappers

import (
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/entities"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/infrastructure/db/schema"
)

func ToUserSchema(entity entities.UserEntity) schema.UserSchema {
	return schema.UserSchema{
		ID:       entity.ID,
		Name:     entity.Name,
		Email:    entity.Email,
		Username: entity.Username,
		Password: entity.Password,
	}
}

func FromUserSchema(schema schema.UserSchema) entities.UserEntity {
	return entities.UserEntity{
		ID:       schema.ID,
		Name:     schema.Name,
		Email:    schema.Email,
		Username: schema.Username,
		Password: schema.Password,
	}
}
