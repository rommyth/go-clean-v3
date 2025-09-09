package gorm

import (
	"go-clean-v3/internal/domain/user"
	"go-clean-v3/internal/infrastructure/persistence/gorm/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// toUserModel converts domain User to GORM model
func toUserModel(u *user.User) *models.UserModel {
	return &models.UserModel{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

// toUserDomain converts GORM model to domain User
func toUserDomain(m *models.UserModel) *user.User {
	return &user.User{
		ID:       m.ID,
		Name:     m.Name,
		Email:    m.Email,
		Password: m.Password,
	}
}

// Create implements user.UserRepositoryInterface.
func (u *userRepository) Create(user *user.User) error {
	return u.db.Create(toUserModel(user)).Error
}

// Delete implements user.UserRepositoryInterface.
func (u *userRepository) Delete(id int64) error {
	return u.db.Delete(&models.UserModel{}, "id = ?", id).Error
}

// GetByEmail implements user.UserRepositoryInterface.
func (u *userRepository) GetByEmail(email string) (*user.User, error) {
	var model models.UserModel
	if err := u.db.Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return toUserDomain(&model), nil
}

// GetByID implements user.UserRepositoryInterface.
func (u *userRepository) GetByID(id int64) (*user.User, error) {
	var model models.UserModel
	if err := u.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return toUserDomain(&model), nil
}	

// Update implements user.UserRepositoryInterface.
func (u *userRepository) Update(user *user.User) error {
	return u.db.Save(toUserModel(user)).Error
}

func NewUserRepository(db *gorm.DB) user.UserRepositoryInterface {
	return &userRepository{db: db}
}
