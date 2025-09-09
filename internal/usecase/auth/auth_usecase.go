package auth

import (
	"context"
	"go-clean-v3/internal/domain/auth"
	"go-clean-v3/internal/domain/user"
	userReq "go-clean-v3/internal/usecase/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo user.UserRepositoryInterface
	authService auth.AuthServiceInterface
}

func NewAuthUsecase(userRepo user.UserRepositoryInterface, authService auth.AuthServiceInterface) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		authService: authService,
	}
}

func (a *AuthUsecase) GetUserFromContext(ctx context.Context) (*user.User, error) {
	return nil, nil
}

func (a *AuthUsecase) Login(ctx context.Context, req userReq.LoginUserRequest) (string, error) {
	dbUser, err := a.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	token, err := a.authService.GenerateToken(dbUser)
	if err != nil {
		return "", err
	}

	return token, nil
}