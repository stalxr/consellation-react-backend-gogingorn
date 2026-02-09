package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims - структура для JWT токена
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// InitJWT инициализирует секретный ключ для подписи JWT токенов
// Читает ключ из переменной окружения JWT_SECRET
// Если ключ не задан, использует дефолтное значение (только для разработки!)
func InitJWT() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production" // Дефолтный ключ для разработки
	}
	jwtSecret = []byte(secret)
}

// GenerateToken генерирует JWT access токен для пользователя
// Токен содержит информацию о пользователе (ID, email, роль) и действителен 24 часа
// Возвращает строку с токеном или ошибку при неудаче
func GenerateToken(userID uuid.UUID, email, role string) (string, error) {
	// Время жизни токена - 24 часа
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:   jwt.NewNumericDate(time.Now()),
			NotBefore:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken генерирует JWT refresh токен для обновления access токена
// Refresh токен действителен 30 дней и используется для получения нового access токена
// Возвращает строку с токеном или ошибку при неудаче
func GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:   jwt.NewNumericDate(time.Now()),
			NotBefore:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken проверяет валидность JWT токена и извлекает из него данные пользователя
// Проверяет подпись токена, срок действия и формат
// Возвращает структуру с данными пользователя (claims) или ошибку при невалидном токене
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи токена")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("невалидный токен")
}

