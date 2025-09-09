package handler

import (
	"go-clean-v3/internal/infrastructure/delivery/http/middleware"
	"go-clean-v3/internal/usecase/user"
	"go-clean-v3/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UserHandler struct {
	userUsecase *user.UserUsecase
}

func NewUserHandler(userUsecase *user.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// REgister handles user registration
func (h *UserHandler) Register(c echo.Context) error {
	var req user.RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-Register-1] Bind error: %v", err)
		return err
	}

	// Call usecase
	userResp, err := h.userUsecase.Register(c.Request().Context(), req)
	if err != nil {
		log.Errorf("[UserHandler-Register-2] Usecase error: %v", err)
		return err
	}

	return response.JSON(c, http.StatusCreated, userResp)
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return err
	}

	userResp, err := h.userUsecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "User not found", err)
	}

	return response.JSON(c, http.StatusOK, userResp)
}