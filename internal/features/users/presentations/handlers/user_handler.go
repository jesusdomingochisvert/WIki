package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/application/usecases"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/presentations/dto"
)

type UserHandler struct {
	getAllUsersUseCase usecases.GetAllUsersUseCase
}

func NewUserHandler(getAllUsersUseCase usecases.GetAllUsersUseCase) *UserHandler {
	return &UserHandler{getAllUsersUseCase: getAllUsersUseCase}
}

func (h *UserHandler) GetAllUsers(ctx *gin.Context) ([]dto.UserDto, error) {
	users, err := h.getAllUsersUseCase.Execute(ctx.Request.Context())
	if err != nil {
		return nil, err
	}

	return users, nil
}
