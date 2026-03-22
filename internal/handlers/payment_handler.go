package handlers

import (
	"net/http"

	"charity-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentHandler - обработчик для webhook'ов платежной системы
type PaymentHandler struct {
	donationService *services.DonationService
}

// NewPaymentHandler - создание нового обработчика платежей
func NewPaymentHandler(donationService *services.DonationService) *PaymentHandler {
	return &PaymentHandler{donationService: donationService}
}

// PaymentCallback - callback от платежной системы (CloudPayments и т.д.)
// @Summary Callback от платежной системы
// @Description Принимает уведомления об успешных платежах от платежной системы
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body PaymentCallbackRequest true "Данные платежа"
// @Success 200 {object} map[string]string
// @Router /payments/callback [post]
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	var req PaymentCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим UUID доната
	donationID, err := uuid.Parse(req.DonationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID доната"})
		return
	}

	// Проверяем статус платежа
	if req.Status != "completed" && req.Status != "success" {
		c.JSON(http.StatusOK, gin.H{"message": "платеж не завершен", "status": req.Status})
		return
	}

	// Завершаем пожертвование
	if err := h.donationService.CompleteDonation(donationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "платеж успешно обработан"})
}

// PaymentCallbackRequest - структура запроса от платежной системы
type PaymentCallbackRequest struct {
	DonationID string  `json:"donation_id" binding:"required"`
	Status     string  `json:"status" binding:"required"` // completed, pending, failed
	Amount     float64 `json:"amount"`
	PaymentID  string  `json:"payment_id,omitempty"` // ID транзакции в платежной системе
}

// ManualCompleteDonation - ручное подтверждение пожертвования (для админа)
// @Summary Ручное подтверждение пожертвования
// @Description Позволяет админу вручную отметить пожертвование как завершенное
// @Tags Payment (Admin)
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID пожертвования"
// @Success 200 {object} map[string]string
// @Router /donations/{id}/complete [post]
func (h *PaymentHandler) ManualCompleteDonation(c *gin.Context) {
	donationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	if err := h.donationService.CompleteDonation(donationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "пожертвование успешно завершено"})
}
