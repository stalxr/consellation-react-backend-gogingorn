package services

import (
	"errors"
	"strings"

	"charity-backend/internal/models"
	"charity-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DreamService - сервис для работы с мечтами
type DreamService struct {
	dreamRepo *repository.DreamRepository
}

// NewDreamService - создание нового сервиса мечтаний
func NewDreamService(dreamRepo *repository.DreamRepository) *DreamService {
	return &DreamService{dreamRepo: dreamRepo}
}

// CreateDreamRequest - запрос на создание мечты
type CreateDreamRequest struct {
	Title           string  `json:"title" binding:"required"`
	ShortDescription *string `json:"short_description,omitempty"`
	FullDescription *string `json:"full_description,omitempty"`
	TargetAmount    float64 `json:"target_amount" binding:"gte=0"`
	CoverImage      *string `json:"cover_image,omitempty"`
}

// DreamsListResponse - ответ со списком мечтаний
type DreamsListResponse struct {
	Data  []models.Dream `json:"data"`
	Total int64          `json:"total"`
}

// Create - создание новой мечты (только для админа)
func (s *DreamService) Create(req CreateDreamRequest) (*models.Dream, error) {
	// Генерируем slug из title
	slug := generateSlug(req.Title)

	// Проверяем уникальность slug
	existing, _ := s.dreamRepo.GetBySlug(slug)
	if existing != nil {
		// Если slug уже существует, добавляем UUID
		slug = slug + "-" + uuid.New().String()[:8]
	}

	dream := &models.Dream{
		ID:              uuid.New(),
		Title:           req.Title,
		Slug:            slug,
		ShortDescription: req.ShortDescription,
		FullDescription: req.FullDescription,
		TargetAmount:    req.TargetAmount,
		CollectedAmount: 0,
		Status:          "active",
		CoverImage:      req.CoverImage,
		GalleryImages:   []string{},
	}

	if err := s.dreamRepo.Create(dream); err != nil {
		return nil, errors.New("ошибка при создании мечты")
	}

	return dream, nil
}

// GetByID - получение мечты по ID
func (s *DreamService) GetByID(id uuid.UUID) (*models.Dream, error) {
	dream, err := s.dreamRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("мечта не найдена")
		}
		return nil, err
	}
	return dream, nil
}

// GetAll - получение списка мечтаний с фильтрацией и пагинацией
func (s *DreamService) GetAll(status string, page, limit int) (*DreamsListResponse, error) {
	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 9
	}
	if limit > 100 {
		limit = 100
	}

	dreams, total, err := s.dreamRepo.GetAll(status, page, limit)
	if err != nil {
		return nil, errors.New("ошибка при получении списка мечтаний")
	}

	return &DreamsListResponse{
		Data:  dreams,
		Total: total,
	}, nil
}

// Update - обновление мечты
func (s *DreamService) Update(id uuid.UUID, dream *models.Dream) error {
	existing, err := s.dreamRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("мечта не найдена")
		}
		return err
	}

	// Обновляем поля
	existing.Title = dream.Title
	if dream.ShortDescription != nil {
		existing.ShortDescription = dream.ShortDescription
	}
	if dream.FullDescription != nil {
		existing.FullDescription = dream.FullDescription
	}
	if dream.TargetAmount > 0 {
		existing.TargetAmount = dream.TargetAmount
	}
	if dream.Status != "" {
		existing.Status = dream.Status
	}
	if dream.CoverImage != nil {
		existing.CoverImage = dream.CoverImage
	}
	if dream.GalleryImages != nil {
		existing.GalleryImages = dream.GalleryImages
	}

	return s.dreamRepo.Update(existing)
}

// Delete - удаление мечты
func (s *DreamService) Delete(id uuid.UUID) error {
	_, err := s.dreamRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("мечта не найдена")
		}
		return err
	}

	return s.dreamRepo.Delete(id)
}

// generateSlug - генерация slug из строки
func generateSlug(s string) string {
	// Простая генерация slug (можно улучшить)
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "ё", "e")
	s = strings.ReplaceAll(s, "й", "y")
	s = strings.ReplaceAll(s, "ц", "ts")
	s = strings.ReplaceAll(s, "у", "u")
	s = strings.ReplaceAll(s, "к", "k")
	s = strings.ReplaceAll(s, "е", "e")
	s = strings.ReplaceAll(s, "н", "n")
	s = strings.ReplaceAll(s, "г", "g")
	s = strings.ReplaceAll(s, "ш", "sh")
	s = strings.ReplaceAll(s, "щ", "sch")
	s = strings.ReplaceAll(s, "з", "z")
	s = strings.ReplaceAll(s, "х", "h")
	s = strings.ReplaceAll(s, "ъ", "")
	s = strings.ReplaceAll(s, "ф", "f")
	s = strings.ReplaceAll(s, "ы", "y")
	s = strings.ReplaceAll(s, "в", "v")
	s = strings.ReplaceAll(s, "а", "a")
	s = strings.ReplaceAll(s, "п", "p")
	s = strings.ReplaceAll(s, "р", "r")
	s = strings.ReplaceAll(s, "о", "o")
	s = strings.ReplaceAll(s, "л", "l")
	s = strings.ReplaceAll(s, "д", "d")
	s = strings.ReplaceAll(s, "ж", "zh")
	s = strings.ReplaceAll(s, "э", "e")
	s = strings.ReplaceAll(s, "я", "ya")
	s = strings.ReplaceAll(s, "ч", "ch")
	s = strings.ReplaceAll(s, "с", "s")
	s = strings.ReplaceAll(s, "м", "m")
	s = strings.ReplaceAll(s, "и", "i")
	s = strings.ReplaceAll(s, "т", "t")
	s = strings.ReplaceAll(s, "ь", "")
	s = strings.ReplaceAll(s, "б", "b")
	s = strings.ReplaceAll(s, "ю", "yu")
	
	// Удаляем все не-латинские символы и дефисы в начале/конце
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	slug := result.String()
	// Удаляем множественные дефисы
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	
	if slug == "" {
		slug = "dream"
	}
	
	return slug
}

