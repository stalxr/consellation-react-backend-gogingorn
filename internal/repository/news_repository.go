package repository

import (
	"charity-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NewsRepository - репозиторий для работы с новостями
type NewsRepository struct {
	db *gorm.DB
}

// NewNewsRepository - создание нового репозитория новостей
func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

// Create - создание новости
func (r *NewsRepository) Create(news *models.News) error {
	return r.db.Create(news).Error
}

// GetAll - получение всех новостей
func (r *NewsRepository) GetAll() ([]models.News, error) {
	var news []models.News
	err := r.db.Order("published_at DESC, created_at DESC").Find(&news).Error
	if err != nil {
		return nil, err
	}
	return news, nil
}

// GetByID - получение новости по ID
func (r *NewsRepository) GetByID(id uuid.UUID) (*models.News, error) {
	var news models.News
	err := r.db.Where("id = ?", id).First(&news).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

// Update - обновление новости
func (r *NewsRepository) Update(news *models.News) error {
	return r.db.Save(news).Error
}

// Delete - удаление новости
func (r *NewsRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.News{}, id).Error
}

