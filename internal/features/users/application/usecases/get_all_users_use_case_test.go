package usecases

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/entities"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/mocks"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	ctx := gomock.NewController(t)
	defer ctx.Finish()

	mockUserRepository := mocks.NewMockUserRepository(ctx)
	useCase := NewGetAllUsersUseCase(mockUserRepository)

	t.Run("Success - return list of users", func(t *testing.T) {
		expectedUsers := []entities.UserEntity{
			{ID: "1", Name: "test1", Email: "test1@test.com", Username: "test1", Password: "test1"},
			{ID: "2", Name: "test2", Email: "test2@test.com", Username: "test2", Password: "test2"},
		}

		mockUserRepository.EXPECT().GetAllUsers(gomock.Any()).Return(expectedUsers, nil)

		users, err := useCase.Execute(context.Background())
		if err != nil {
			t.Fatalf("No se esperaba error, pero ocurrió: %v", err)
		}

		if len(expectedUsers) != len(users) {
			t.Errorf("Se esperaban %v, pero se obtuvo %v", expectedUsers, users)
		}
	})

	t.Run("Error - repository fails", func(t *testing.T) {
		mockUserRepository.EXPECT().GetAllUsers(gomock.Any()).Return(nil, errors.New("repo error"))

		users, err := useCase.Execute(context.Background())
		if err == nil {
			t.Fatalf("Se esperaba un error, pero no ocurrió ninguno")
		}

		if users != nil {
			t.Errorf("Se esperaba un slice nil, pero se obtuvo %v", users)
		}
	})
}
