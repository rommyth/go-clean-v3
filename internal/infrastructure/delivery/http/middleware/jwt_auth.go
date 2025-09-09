package middleware

import (
	"go-clean-v3/internal/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	ecfg := echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
	}

	return echojwt.WithConfig(ecfg)
}

func GetUserIDFromToken(c echo.Context) (int64, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return -1, echo.ErrUnauthorized
	}
	return int64(userID), nil
}