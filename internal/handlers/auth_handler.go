package handlers

import (
	"net/http"

	"charity-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler - обработчик для аутентификации
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler - создание нового обработчика аутентификации
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register - регистрация донора
// @Summary Регистрация донора (Личный кабинет)
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body services.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} services.AuthTokenResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Login - вход (админ или донор)
// @Summary Вход (Админ или Донор)
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body services.LoginRequest true "Данные для входа"
// @Success 200 {object} services.AuthTokenResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// GetMe - получение профиля текущего пользователя
// @Summary Получить профиль текущего пользователя
// @Tags Auth
// @Security bearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	user, err := h.authService.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

