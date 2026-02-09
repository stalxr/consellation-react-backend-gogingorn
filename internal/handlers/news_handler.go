package handlers

import (
	"net/http"

	"charity-backend/internal/models"
	"charity-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// NewsHandler - обработчик для новостей
type NewsHandler struct {
	newsService *services.NewsService
}

// NewNewsHandler - создание нового обработчика новостей
func NewNewsHandler(newsService *services.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

// GetNews - получение списка новостей
// @Summary Получить список новостей
// @Tags News
// @Produce json
// @Success 200 {array} models.News
// @Router /news [get]
func (h *NewsHandler) GetNews(c *gin.Context) {
	news, err := h.newsService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

// CreateNews - публикация новости (только админ)
// @Summary Опубликовать новость
// @Tags News (Admin)
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param body body models.News true "Данные новости"
// @Success 201 {object} models.News
// @Router /news [post]
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdNews, err := h.newsService.Create(&news)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdNews)
}

