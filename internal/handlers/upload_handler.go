package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadHandler - обработчик для загрузки файлов
type UploadHandler struct{}

// NewUploadHandler - создание нового обработчика загрузки
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage - загрузка картинок (для админки)
// @Summary Загрузка картинок (для админки)
// @Tags Upload
// @Security bearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Изображение"
// @Success 200 {object} map[string]string
// @Router /upload [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл не найден"})
		return
	}

	// Валидация расширения файла (только изображения)
	ext := filepath.Ext(file.Filename)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неподдерживаемый формат файла. Разрешены: jpg, jpeg, png, gif, webp"})
		return
	}

	// Генерируем уникальное имя файла
	filename := uuid.New().String() + ext
	uploadDir := filepath.Join("uploads", "images")
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

	// Возвращаем URL файла
	c.JSON(http.StatusOK, gin.H{
		"url": "/uploads/images/" + filename, // В реальном проекте это будет полный URL
	})
}
