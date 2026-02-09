package services

import (
	"errors"
	"strings"

	"charity-backend/internal/models"
	"charity-backend/internal/repository"
	"charity-backend/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService - сервис для аутентификации
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService - создание нового сервиса аутентификации
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// RegisterRequest - запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginRequest - запрос на вход
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthTokenResponse - ответ с токенами
type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register регистрирует нового пользователя (донора) в системе
// Хэширует пароль с помощью bcrypt, создает запись в БД и генерирует JWT токены
// Возвращает access и refresh токены для немедленной авторизации
func (s *AuthService) Register(req RegisterRequest) (*AuthTokenResponse, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Хэшируем пароль
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("ошибка при хэшировании пароля")
	}

	// Создаем нового пользователя
	user := &models.User{
		ID:           uuid.New(),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         "user", // По умолчанию роль "user"
		TotalDonated: 0,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("ошибка при создании пользователя")
	}

	// Генерируем токены
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("ошибка при генерации токена")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("ошибка при генерации refresh токена")
	}

	return &AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login выполняет аутентификацию пользователя по email и паролю
// Проверяет существование пользователя и корректность пароля
// Возвращает JWT токены (access и refresh) при успешной авторизации
func (s *AuthService) Login(req LoginRequest) (*AuthTokenResponse, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("неверный email или пароль")
		}
		return nil, err
	}

	// Проверяем пароль
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("неверный email или пароль")
	}

	// Генерируем токены
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("ошибка при генерации токена")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("ошибка при генерации refresh токена")
	}

	return &AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetUserByID - получение пользователя по ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

