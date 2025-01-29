package mappers

import (
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/entities"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/presentations/dto"
)

func ToUserEntity(dto dto.UserDto) entities.UserEntity {
	return entities.UserEntity{
		ID:       dto.ID,
		Name:     dto.Name,
		Email:    dto.Email,
		Username: dto.Username,
		Password: dto.Password,
	}
}

func FromUserEntity(entity entities.UserEntity) dto.UserDto {
	return dto.UserDto{
		ID:       entity.ID,
		Name:     entity.Name,
		Email:    entity.Email,
		Username: entity.Username,
		Password: entity.Password,
	}
}
