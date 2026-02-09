package repository

import (
	"charity-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository - репозиторий для работы с пользователями
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository - создание нового репозитория пользователей
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create - создание нового пользователя
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByEmail - получение пользователя по email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID - получение пользователя по ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update - обновление пользователя
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// UpdateTotalDonated - обновление суммы пожертвований пользователя
func (r *UserRepository) UpdateTotalDonated(userID uuid.UUID, amount float64) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("total_donated", gorm.Expr("total_donated + ?", amount)).Error
}

