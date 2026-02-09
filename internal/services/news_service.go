package services

import (
	"errors"
	"time"

	"charity-backend/internal/models"
	"charity-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NewsService - сервис для работы с новостями
type NewsService struct {
	newsRepo *repository.NewsRepository
}

// NewNewsService - создание нового сервиса новостей
func NewNewsService(newsRepo *repository.NewsRepository) *NewsService {
	return &NewsService{newsRepo: newsRepo}
}

// GetAll - получение всех новостей
func (s *NewsService) GetAll() ([]models.News, error) {
	return s.newsRepo.GetAll()
}

// Create - создание новости
func (s *NewsService) Create(news *models.News) (*models.News, error) {
	news.ID = uuid.New()
	now := time.Now()
	news.PublishedAt = &now

	if err := s.newsRepo.Create(news); err != nil {
		return nil, errors.New("ошибка при создании новости")
	}

	return news, nil
}

// GetByID - получение новости по ID
func (s *NewsService) GetByID(id uuid.UUID) (*models.News, error) {
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("новость не найдена")
		}
		return nil, err
	}
	return news, nil
}

// Update - обновление новости
func (s *NewsService) Update(news *models.News) error {
	_, err := s.newsRepo.GetByID(news.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("новость не найдена")
		}
		return err
	}

	return s.newsRepo.Update(news)
}

// Delete - удаление новости
func (s *NewsService) Delete(id uuid.UUID) error {
	_, err := s.newsRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("новость не найдена")
		}
		return err
	}

	return s.newsRepo.Delete(id)
}

