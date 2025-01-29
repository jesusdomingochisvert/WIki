package usecases

import (
	"context"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/application/mappers"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/repositories"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/presentations/dto"
)

type GetAllUsersUseCase interface {
	Execute(ctx context.Context) ([]dto.UserDto, error)
}

type getAllUsersUseCase struct {
	userRepository repositories.UserRepository
}

func NewGetAllUsersUseCase(userRepository repositories.UserRepository) GetAllUsersUseCase {
	return &getAllUsersUseCase{userRepository: userRepository}
}

func (r *getAllUsersUseCase) Execute(ctx context.Context) ([]dto.UserDto, error) {
	users, err := r.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userDtos []dto.UserDto
	for _, user := range users {
		userDto := mappers.FromUserEntity(user)
		userDtos = append(userDtos, userDto)
	}

	return userDtos, nil
}
