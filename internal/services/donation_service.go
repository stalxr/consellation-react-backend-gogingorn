package services

import (
	"errors"
	"fmt"

	"charity-backend/internal/models"
	"charity-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DonationService - сервис для работы с пожертвованиями
type DonationService struct {
	donationRepo *repository.DonationRepository
	dreamRepo    *repository.DreamRepository
	userRepo     *repository.UserRepository
}

// NewDonationService - создание нового сервиса пожертвований
func NewDonationService(
	donationRepo *repository.DonationRepository,
	dreamRepo *repository.DreamRepository,
	userRepo *repository.UserRepository,
) *DonationService {
	return &DonationService{
		donationRepo: donationRepo,
		dreamRepo:    dreamRepo,
		userRepo:     userRepo,
	}
}

// CreateDonationRequest - запрос на создание пожертвования
type CreateDonationRequest struct {
	Amount      float64    `json:"amount" binding:"required,gt=0"`
	DreamID     uuid.UUID  `json:"dream_id" binding:"required"`
	UserID      *uuid.UUID `json:"user_id,omitempty"` // Может быть nil для анонимных
	Email       *string    `json:"email,omitempty"`
	IsAnonymous bool       `json:"is_anonymous,omitempty"`
	Comment     *string    `json:"comment,omitempty"`
}

// PaymentResponse - ответ с URL для оплаты
type PaymentResponse struct {
	PaymentURL string `json:"payment_url"`
}

// CreateDonation - создание пожертвования
func (s *DonationService) CreateDonation(req CreateDonationRequest) (*PaymentResponse, error) {
	// Проверяем существование мечты
	dream, err := s.dreamRepo.GetByID(req.DreamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("мечта не найдена")
		}
		return nil, err
	}

	// Проверяем, что мечта активна
	if dream.Status != "active" {
		return nil, errors.New("мечта не активна")
	}

	// Создаем пожертвование
	donation := &models.Donation{
		ID:          uuid.New(),
		DreamID:     req.DreamID,
		UserID:      req.UserID,
		Amount:      req.Amount,
		Email:       req.Email,
		IsAnonymous: req.IsAnonymous,
		Comment:     req.Comment,
		Status:      "pending",
	}

	// Генерируем URL для оплаты (в реальном проекте здесь будет интеграция с платежной системой)
	paymentURL := fmt.Sprintf("https://payment.example.com/pay?donation_id=%s&amount=%.2f", donation.ID, donation.Amount)
	donation.PaymentURL = &paymentURL

	if err := s.donationRepo.Create(donation); err != nil {
		return nil, errors.New("ошибка при создании пожертвования")
	}

	return &PaymentResponse{
		PaymentURL: paymentURL,
	}, nil
}

// GetByUserID - получение всех пожертвований пользователя
func (s *DonationService) GetByUserID(userID uuid.UUID) ([]models.Donation, error) {
	return s.donationRepo.GetByUserID(userID)
}

// GetByDreamID - получение публичного списка пожертвований для мечты
func (s *DonationService) GetByDreamID(dreamID uuid.UUID) ([]models.Donation, error) {
	return s.donationRepo.GetByDreamID(dreamID)
}

// CompleteDonation - завершение пожертвования (вызывается после успешной оплаты)
func (s *DonationService) CompleteDonation(donationID uuid.UUID) error {
	donation, err := s.donationRepo.GetByID(donationID)
	if err != nil {
		return err
	}

	if donation.Status == "completed" {
		return nil // Уже завершено
	}

	// Обновляем статус
	donation.Status = "completed"
	if err := s.donationRepo.Update(donation); err != nil {
		return err
	}

	// Обновляем собранную сумму мечты
	if err := s.dreamRepo.UpdateCollectedAmount(donation.DreamID, donation.Amount); err != nil {
		return err
	}

	// Обновляем сумму пожертвований пользователя (если не анонимное)
	if donation.UserID != nil {
		if err := s.userRepo.UpdateTotalDonated(*donation.UserID, donation.Amount); err != nil {
			return err
		}
	}

	return nil
}

