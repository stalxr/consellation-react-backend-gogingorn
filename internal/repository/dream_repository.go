package repository

import (
	"charity-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DreamRepository - репозиторий для работы с мечтами
type DreamRepository struct {
	db *gorm.DB
}

// NewDreamRepository - создание нового репозитория мечтаний
func NewDreamRepository(db *gorm.DB) *DreamRepository {
	return &DreamRepository{db: db}
}

// Create - создание новой мечты
func (r *DreamRepository) Create(dream *models.Dream) error {
	return r.db.Create(dream).Error
}

// GetByID - получение мечты по ID
func (r *DreamRepository) GetByID(id uuid.UUID) (*models.Dream, error) {
	var dream models.Dream
	err := r.db.Where("id = ?", id).First(&dream).Error
	if err != nil {
		return nil, err
	}
	return &dream, nil
}

// GetBySlug - получение мечты по slug
func (r *DreamRepository) GetBySlug(slug string) (*models.Dream, error) {
	var dream models.Dream
	err := r.db.Where("slug = ?", slug).First(&dream).Error
	if err != nil {
		return nil, err
	}
	return &dream, nil
}

// GetAll - получение списка мечтаний с фильтрацией и пагинацией
func (r *DreamRepository) GetAll(status string, page, limit int) ([]models.Dream, int64, error) {
	var dreams []models.Dream
	var total int64

	query := r.db.Model(&models.Dream{})

	// Фильтр по статусу
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Подсчет общего количества
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Пагинация
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&dreams).Error; err != nil {
		return nil, 0, err
	}

	return dreams, total, nil
}

// Update - обновление мечты
func (r *DreamRepository) Update(dream *models.Dream) error {
	return r.db.Save(dream).Error
}

// Delete - удаление мечты (мягкое удаление)
func (r *DreamRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Dream{}, id).Error
}

// UpdateCollectedAmount - обновление собранной суммы
func (r *DreamRepository) UpdateCollectedAmount(dreamID uuid.UUID, amount float64) error {
	return r.db.Model(&models.Dream{}).
		Where("id = ?", dreamID).
		Update("collected_amount", gorm.Expr("collected_amount + ?", amount)).Error
}

