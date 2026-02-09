package services

import (
	"errors"

	"charity-backend/internal/models"
	"charity-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportService - сервис для работы с отчетами
type ReportService struct {
	reportRepo *repository.ReportRepository
}

// NewReportService - создание нового сервиса отчетов
func NewReportService(reportRepo *repository.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

// GetAll - получение всех отчетов
func (s *ReportService) GetAll() ([]models.Report, error) {
	return s.reportRepo.GetAll()
}

// Create - создание отчета
func (s *ReportService) Create(report *models.Report) (*models.Report, error) {
	// Валидация месяца
	if report.Month < 1 || report.Month > 12 {
		return nil, errors.New("месяц должен быть от 1 до 12")
	}

	report.ID = uuid.New()

	if err := s.reportRepo.Create(report); err != nil {
		return nil, errors.New("ошибка при создании отчета")
	}

	return report, nil
}

// GetByID - получение отчета по ID
func (s *ReportService) GetByID(id uuid.UUID) (*models.Report, error) {
	report, err := s.reportRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("отчет не найден")
		}
		return nil, err
	}
	return report, nil
}

// Update - обновление отчета
func (s *ReportService) Update(report *models.Report) error {
	_, err := s.reportRepo.GetByID(report.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("отчет не найден")
		}
		return err
	}

	return s.reportRepo.Update(report)
}

// Delete - удаление отчета
func (s *ReportService) Delete(id uuid.UUID) error {
	_, err := s.reportRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("отчет не найден")
		}
		return err
	}

	return s.reportRepo.Delete(id)
}

