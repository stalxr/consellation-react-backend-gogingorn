package handlers

import (
	"net/http"
	"strconv"

	"charity-backend/internal/models"
	"charity-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DreamHandler - обработчик для мечтаний
type DreamHandler struct {
	dreamService *services.DreamService
}

// NewDreamHandler - создание нового обработчика мечтаний
func NewDreamHandler(dreamService *services.DreamService) *DreamHandler {
	return &DreamHandler{dreamService: dreamService}
}

// GetDreams - получение списка всех мечтаний (с фильтрами)
// @Summary Список всех мечтаний (с фильтрами)
// @Tags Dreams (Public)
// @Produce json
// @Param status query string false "Статус (active, completed)"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(9)
// @Success 200 {object} services.DreamsListResponse
// @Router /dreams [get]
func (h *DreamHandler) GetDreams(c *gin.Context) {
	status := c.Query("status")
	
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "9"))
	if err != nil || limit < 1 {
		limit = 9
	}

	result, err := h.dreamService.GetAll(status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetDreamByID - получение детальной информации о мечте
// @Summary Детальная информация о мечте
// @Tags Dreams (Public)
// @Produce json
// @Param id path string true "ID мечты"
// @Success 200 {object} models.Dream
// @Router /dreams/{id} [get]
func (h *DreamHandler) GetDreamByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	dream, err := h.dreamService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dream)
}

// CreateDream - создание новой карточки мечты (только админ)
// @Summary Создать новую карточку мечты
// @Tags Dreams (Admin)
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param body body services.CreateDreamRequest true "Данные мечты"
// @Success 201 {object} models.Dream
// @Router /dreams [post]
func (h *DreamHandler) CreateDream(c *gin.Context) {
	var req services.CreateDreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dream, err := h.dreamService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dream)
}

// UpdateDream - редактирование мечты (только админ)
// @Summary Редактировать мечту
// @Tags Dreams (Admin)
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID мечты"
// @Param body body models.Dream true "Обновленные данные мечты"
// @Success 200 {object} models.Dream
// @Router /dreams/{id} [put]
func (h *DreamHandler) UpdateDream(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	var dream models.Dream
	if err := c.ShouldBindJSON(&dream); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.dreamService.Update(id, &dream); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedDream, err := h.dreamService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedDream)
}

// DeleteDream - удаление мечты (только админ)
// @Summary Удалить мечту
// @Tags Dreams (Admin)
// @Security bearerAuth
// @Param id path string true "ID мечты"
// @Success 200 {object} gin.H
// @Router /dreams/{id} [delete]
func (h *DreamHandler) DeleteDream(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	if err := h.dreamService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "мечта успешно удалена"})
}

