package user

import (
	"context"
	"go-clean-v3/internal/domain/auth"
	"go-clean-v3/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo    user.UserRepositoryInterface
	authService auth.AuthServiceInterface
}

func NewUserUsecase(userRepo user.UserRepositoryInterface, authService auth.AuthServiceInterface) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (u *UserUsecase) Register(ctx context.Context, req RegisterUserRequest) (*UserResponse, error) {
	if _, err := u.userRepo.GetByEmail(req.Email); err == nil {
		return nil, err
	}

	// hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user domain model
	newUser := &user.User{
		Name: req.Name,
		Email: req.Email,
		Password: string(hashPassword),
	}

	// save to repository
	if err := u.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	// Return Response DTO
	return &UserResponse{
		ID: newUser.ID,
		Name: newUser.Name,
		Email: newUser.Email,
	}, nil
}

func (u *UserUsecase) Login(ctx context.Context, req LoginUserRequest) (string, error) {
	existUser, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}

	// Compoare Password
	if err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(existUser.Password)); err != nil {
		return "", user.ErrUserNotFound
	}

	// Generate JWT token
	token, err := u.authService.GenerateToken(existUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) GetProfile(ctx context.Context, userID int64) (*UserResponse, error) {
	userData, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID: userData.ID,
		Name: userData.Name,
		Email: userData.Email,
	}, nil
}