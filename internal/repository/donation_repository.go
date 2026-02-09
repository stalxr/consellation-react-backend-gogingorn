package repository

import (
	"charity-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DonationRepository - репозиторий для работы с пожертвованиями
type DonationRepository struct {
	db *gorm.DB
}

// NewDonationRepository - создание нового репозитория пожертвований
func NewDonationRepository(db *gorm.DB) *DonationRepository {
	return &DonationRepository{db: db}
}

// Create - создание нового пожертвования
func (r *DonationRepository) Create(donation *models.Donation) error {
	return r.db.Create(donation).Error
}

// GetByID - получение пожертвования по ID
func (r *DonationRepository) GetByID(id uuid.UUID) (*models.Donation, error) {
	var donation models.Donation
	err := r.db.Where("id = ?", id).First(&donation).Error
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

// GetByUserID - получение всех пожертвований пользователя
func (r *DonationRepository) GetByUserID(userID uuid.UUID) ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Where("user_id = ?", userID).
		Preload("Dream").
		Order("created_at DESC").
		Find(&donations).Error
	if err != nil {
		return nil, err
	}
	return donations, nil
}

// GetByDreamID - получение всех пожертвований для мечты (публичный список)
func (r *DonationRepository) GetByDreamID(dreamID uuid.UUID) ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Where("dream_id = ? AND status = ?", dreamID, "completed").
		Preload("User").
		Order("created_at DESC").
		Find(&donations).Error
	if err != nil {
		return nil, err
	}
	return donations, nil
}

// Update - обновление пожертвования
func (r *DonationRepository) Update(donation *models.Donation) error {
	return r.db.Save(donation).Error
}

