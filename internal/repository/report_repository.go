package repository

import (
	"charity-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepository - репозиторий для работы с отчетами
type ReportRepository struct {
	db *gorm.DB
}

// NewReportRepository - создание нового репозитория отчетов
func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// Create - создание отчета
func (r *ReportRepository) Create(report *models.Report) error {
	return r.db.Create(report).Error
}

// GetAll - получение всех отчетов
func (r *ReportRepository) GetAll() ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Order("year DESC, month DESC").Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}

// GetByID - получение отчета по ID
func (r *ReportRepository) GetByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// Update - обновление отчета
func (r *ReportRepository) Update(report *models.Report) error {
	return r.db.Save(report).Error
}

// Delete - удаление отчета
func (r *ReportRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Report{}, id).Error
}

