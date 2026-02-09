package handlers

import (
	"net/http"

	"charity-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DonationHandler - обработчик для пожертвований
type DonationHandler struct {
	donationService *services.DonationService
}

// NewDonationHandler - создание нового обработчика пожертвований
func NewDonationHandler(donationService *services.DonationService) *DonationHandler {
	return &DonationHandler{donationService: donationService}
}

// CreatePayment - создание платежа
// @Summary Создать платеж
// @Description Если user_id не передан, считается анонимным. Возвращает ссылку на эквайринг.
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body services.CreateDonationRequest true "Данные платежа"
// @Success 200 {object} services.PaymentResponse
// @Router /donations/pay [post]
func (h *DonationHandler) CreatePayment(c *gin.Context) {
	var req services.CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Если пользователь авторизован, используем его ID
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uuid.UUID)
		req.UserID = &uid
	}

	response, err := h.donationService.CreateDonation(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMyDonations - история моих пожертвований
// @Summary История моих пожертвований
// @Tags User Cabinet
// @Security bearerAuth
// @Produce json
// @Success 200 {array} object
// @Router /donations/my [get]
func (h *DonationHandler) GetMyDonations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	donations, err := h.donationService.GetByUserID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формируем ответ согласно OpenAPI
	result := make([]map[string]interface{}, 0)
	for _, donation := range donations {
		dreamTitle := ""
		if donation.Dream.ID != uuid.Nil {
			dreamTitle = donation.Dream.Title
		}
		
		item := map[string]interface{}{
			"id":         donation.ID.String(),
			"amount":     donation.Amount,
			"dream_title": dreamTitle,
			"date":       donation.CreatedAt,
			"status":     donation.Status,
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, result)
}

// GetDreamDonors - список тех, кто помог этой мечте
// @Summary Список тех, кто помог этой мечте
// @Tags Dreams (Public)
// @Produce json
// @Param id path string true "ID мечты"
// @Success 200 {array} object
// @Router /dreams/{id}/donors [get]
func (h *DonationHandler) GetDreamDonors(c *gin.Context) {
	dreamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	donations, err := h.donationService.GetByDreamID(dreamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формируем публичный ответ
	result := make([]map[string]interface{}, 0)
	for _, donation := range donations {
		donorName := "Анонимный донор"
		if !donation.IsAnonymous && donation.UserID != nil && donation.User != nil {
			donorName = donation.User.FullName
		}

		item := map[string]interface{}{
			"amount":       donation.Amount,
			"donor_name":   donorName,
			"created_at":   donation.CreatedAt,
			"is_anonymous": donation.IsAnonymous,
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, result)
}

