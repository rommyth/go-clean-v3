package auth

import "go-clean-v3/internal/domain/user"

type AuthServiceInterface interface {
	GenerateToken(u *user.User) (string, error)
}