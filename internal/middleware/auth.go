package middleware

import (
	"net/http"
	"strings"

	"charity-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - middleware для проверки JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "отсутствует токен авторизации"})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		token := parts[1]

		// Валидируем токен
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидный токен"})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контексте
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware - middleware для опциональной проверки JWT токена
// Если заголовок Authorization отсутствует, запрос считается анонимным
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидный токен"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware - middleware для проверки прав администратора
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Сначала проверяем авторизацию
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}

		// Проверяем роль администратора
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав"})
			c.Abort()
			return
		}

		c.Next()
	}
}

