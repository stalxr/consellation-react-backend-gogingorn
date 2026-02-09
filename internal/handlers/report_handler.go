package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"charity-backend/internal/models"
	"charity-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReportHandler - обработчик для отчетов
type ReportHandler struct {
	reportService *services.ReportService
}

// NewReportHandler - создание нового обработчика отчетов
func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetReports - получение финансовых отчетов
// @Summary Финансовые отчеты
// @Tags Reports
// @Produce json
// @Success 200 {array} models.Report
// @Router /reports [get]
func (h *ReportHandler) GetReports(c *gin.Context) {
	reports, err := h.reportService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// UploadReport - загрузка отчета (PDF) (только админ)
// @Summary Загрузить отчет (PDF)
// @Tags Reports (Admin)
// @Security bearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "PDF файл"
// @Param title formData string true "Название отчета"
// @Param year formData int true "Год"
// @Success 201 {object} models.Report
// @Router /reports [post]
func (h *ReportHandler) UploadReport(c *gin.Context) {
	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл не найден"})
		return
	}

	// Получаем остальные данные
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "название отчета обязательно"})
		return
	}

	yearStr := c.PostForm("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "год обязателен"})
		return
	}

	// Валидация расширения файла
	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл должен быть в формате PDF"})
		return
	}

	// Сохраняем файл (в реальном проекте здесь будет загрузка в S3 или другое хранилище)
	// Пока просто сохраняем путь
	filename := uuid.New().String() + ext
	uploadDir := filepath.Join("uploads", "reports")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании директории загрузок"})
		return
	}
	filePath := filepath.Join(uploadDir, filename)
	// Сохраняем файл
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении файла"})
		return
	}

	// Получаем месяц из текущей даты (можно добавить в форму)
	now := time.Now()
	month := int(now.Month())

	// Создаем отчет
	report := &models.Report{
		Title:   title,
		Year:    parseYear(yearStr),
		Month:   month,
		FileURL: "/uploads/reports/" + filename, // В реальном проекте это будет полный URL
	}

	createdReport, err := h.reportService.Create(report)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdReport)
}

// parseYear - парсинг года из строки
func parseYear(s string) int {
	year, err := strconv.Atoi(s)
	if err != nil {
		return time.Now().Year()
	}
	return year
}

