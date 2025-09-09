package handler

import (
	"go-clean-v3/internal/usecase/auth"
	"go-clean-v3/internal/usecase/user"
	"go-clean-v3/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase *auth.AuthUsecase
}

func NewAuthHandler(authUsecase *auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req user.LoginUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	token, err := h.authUsecase.Login(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, map[string]interface{}{
		"token": token,
	})
}