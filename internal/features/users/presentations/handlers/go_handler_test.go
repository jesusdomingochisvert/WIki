package handlers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/mocks"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/presentations/dto"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockGetAllUsersUseCase(ctrl)
	handler := NewUserHandler(mockUseCase)

	gin.SetMode(gin.TestMode)

	t.Run("Success - Return all users", func(t *testing.T) {
		mockUsers := []dto.UserDto{
			{ID: "1", Name: "test1", Email: "test1@test.com", Username: "test1", Password: "test1"},
			{ID: "2", Name: "test2", Email: "test2@test.com", Username: "test2", Password: "test2"},
		}

		mockUseCase.EXPECT().Execute(gomock.Any()).Return(mockUsers, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		users, err := handler.GetAllUsers(ctx)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers, users)
	})

	t.Run("Failed - Return error", func(t *testing.T) {
		expectedError := errors.New("test error")

		mockUseCase.EXPECT().Execute(gomock.Any()).Return(nil, expectedError)

		req := httptest.NewRequest(http.MethodGet, "/users", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		users, err := handler.GetAllUsers(ctx)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, expectedError, err)
	})
}
