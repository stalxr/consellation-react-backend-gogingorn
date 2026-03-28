package main

import (
	"log"
	"os"
	"path/filepath"

	"charity-backend/internal/handlers"
	"charity-backend/internal/middleware"
	"charity-backend/internal/models"
	"charity-backend/internal/repository"
	"charity-backend/internal/services"
	"charity-backend/pkg/db"
	"charity-backend/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Загружаем переменные окружения из .env файла
	projectRoot, _ := filepath.Abs("../../")
	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  Файл .env не найден по пути %s, используются переменные окружения системы", envPath)
	}

	// Инициализируем JWT
	utils.InitJWT()

	// Подключаемся к базе данных
	db.Connect()

	// Получаем подключение к БД
	database := db.DB

	// SeedData заполняет базу тестовыми данными (дети с онкологией)
	SeedData(database)

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(database)
	dreamRepo := repository.NewDreamRepository(database)
	donationRepo := repository.NewDonationRepository(database)
	newsRepo := repository.NewNewsRepository(database)
	reportRepo := repository.NewReportRepository(database)

	// Инициализируем сервисы
	authService := services.NewAuthService(userRepo)
	dreamService := services.NewDreamService(dreamRepo)
	donationService := services.NewDonationService(donationRepo, dreamRepo, userRepo)
	newsService := services.NewNewsService(newsRepo)
	reportService := services.NewReportService(reportRepo)

	// Инициализируем обработчики
	authHandler := handlers.NewAuthHandler(authService)
	dreamHandler := handlers.NewDreamHandler(dreamService)
	donationHandler := handlers.NewDonationHandler(donationService)
	newsHandler := handlers.NewNewsHandler(newsService)
	reportHandler := handlers.NewReportHandler(reportService)
	uploadHandler := handlers.NewUploadHandler()
	paymentHandler := handlers.NewPaymentHandler(donationService)

	// Настраиваем Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS (критично для React): разрешаем запросы со всех источников
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"Cache-Control",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Подготовка локальных директорий для загрузок
	if err := os.MkdirAll(filepath.Join("uploads", "images"), 0o755); err != nil {
		log.Fatalf("❌ Ошибка создания директории загрузок: %v", err)
	}
	if err := os.MkdirAll(filepath.Join("uploads", "reports"), 0o755); err != nil {
		log.Fatalf("❌ Ошибка создания директории загрузок: %v", err)
	}

	// Swagger UI endpoint - фронтенд смотрит документацию API
	router.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `<!DOCTYPE html>
<html>
<head>
	<title>Charity API - Swagger UI</title>
	<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
	<script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
	<script>
		window.onload = function() {
			SwaggerUIBundle({
				url: '/openapi.yaml',
				dom_id: '#swagger-ui',
				presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
				layout: "StandaloneLayout"
			});
		};
	</script>
</body>
</html>`)
	})

	// Раздаём openapi.yaml статически
	router.StaticFile("/openapi.yaml", "./openapi.yaml")

	// Настраиваем статическую раздачу файлов из папки uploads
	// Это позволяет открывать загруженные изображения и PDF через браузер
	router.Static("/uploads", "./uploads")

	// API v1 группа
	api := router.Group("/api/v1")
	{
		// Health check endpoint (для Railway)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		// Публичные роуты (Auth)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetMe)
		}

		// Публичные роуты (Dreams)
		dreams := api.Group("/dreams")
		{
			dreams.GET("", dreamHandler.GetDreams)
			dreams.GET("/:id", dreamHandler.GetDreamByID)
			dreams.GET("/:id/donors", donationHandler.GetDreamDonors)

			// Админские роуты (Dreams)
			dreamsAdmin := dreams.Group("", middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				dreamsAdmin.POST("", dreamHandler.CreateDream)
				dreamsAdmin.PUT("/:id", dreamHandler.UpdateDream)
				dreamsAdmin.DELETE("/:id", dreamHandler.DeleteDream)
			}
		}

		// Публичные роуты (Donations)
		donations := api.Group("/donations")
		{
			donations.POST("/pay", middleware.OptionalAuthMiddleware(), donationHandler.CreatePayment)
			donations.GET("/my", middleware.AuthMiddleware(), donationHandler.GetMyDonations)

			// Админские роуты для ручного подтверждения
			donations.POST("/:id/complete", middleware.AuthMiddleware(), middleware.AdminMiddleware(), paymentHandler.ManualCompleteDonation)
		}

		// Webhook для платежной системы (публичный, но с защитой по секрету или IP)
		api.POST("/payments/callback", paymentHandler.PaymentCallback)

		// Публичные роуты (News)
		news := api.Group("/news")
		{
			news.GET("", newsHandler.GetNews)

			// Админские роуты (News)
			newsAdmin := news.Group("", middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				newsAdmin.POST("", newsHandler.CreateNews)
			}
		}

		// Публичные роуты (Reports)
		reports := api.Group("/reports")
		{
			reports.GET("", reportHandler.GetReports)

			// Админские роуты (Reports)
			reportsAdmin := reports.Group("", middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				reportsAdmin.POST("", reportHandler.UploadReport)
			}
		}

		// Админские роуты (Upload)
		upload := api.Group("/upload", middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			upload.POST("", uploadHandler.UploadImage)
		}
	}

	// Получаем порт из переменных окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Сервер запущен на порту %s", port)
	log.Printf("📖 API доступно по адресу: http://localhost:%s/api/v1", port)

	// Запускаем сервер
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Ошибка запуска сервера: %v", err)
	}
}

// SeedData заполняет базу тестовыми данными (дети с онкологией)
func SeedData(database *gorm.DB) {
	log.Println("🌱 Начинаем заполнение тестовыми данными...")

	var count int64
	database.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("ℹ️ Данные уже существуют, пропускаем seed")
		return
	}

	// Создаем админа
	adminPassword, _ := utils.HashPassword("admin123")
	admin := &models.User{
		ID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Email:        "admin@nastenka.ru",
		PasswordHash: adminPassword,
		FullName:     "Администратор Настенька",
		Role:         "admin",
		TotalDonated: 0,
	}
	database.Create(admin)
	log.Println("✅ Админ: admin@nastenka.ru / admin123")

	// Создаем тестового пользователя
	userPassword, _ := utils.HashPassword("user123")
	user := &models.User{
		ID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Email:        "user@example.com",
		PasswordHash: userPassword,
		FullName:     "Иван Петров",
		Role:         "user",
		TotalDonated: 15000,
	}
	database.Create(user)
	log.Println("✅ Юзер: user@example.com / user123")

	// Дети с онкологией
	dreams := []models.Dream{
		{
			ID:               uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Title:            "Лечение лейкемии для Маши, 7 лет",
			Slug:             "lechenie-leykemii-mashi-7-let",
			ShortDescription: strPtr("Маше нужен курс химиотерапии и пересадка костного мозга"),
			FullDescription:  strPtr("<p>У Маши острая лимфобластная лейкемия. Необходим курс химиотерапии с пересадкой костного мозга. Стоимость: 2,800,000 руб.</p>"),
			TargetAmount:     2800000,
			CollectedAmount:  950000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/masha.jpg"),
		},
		{
			ID:               uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Title:            "Операция по удалению опухоли мозга для Вани, 5 лет",
			Slug:             "operaciya-opuhol-mozga-vanya-5-let",
			ShortDescription: strPtr("Сложная нейрохирургическая операция в Израиле"),
			FullDescription:  strPtr("<p>У Вани медуллобластома. Операция в Израиле. Стоимость: 3,500,000 руб.</p>"),
			TargetAmount:     3500000,
			CollectedAmount:  2100000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/vanya.jpg"),
		},
		{
			ID:               uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			Title:            "Протонная терапия для Кати, 12 лет (рак глаза)",
			Slug:             "protonnaya-terapiya-katya-12-let",
			ShortDescription: strPtr("Современное лечение ретинобластомы в Германии"),
			FullDescription:  strPtr("<p>У Кати ретинобластома. Протонная терапия в Гейдельберге. Стоимость: 1,800,000 руб.</p>"),
			TargetAmount:     1800000,
			CollectedAmount:  1800000,
			Status:           "completed",
			CoverImage:       strPtr("/uploads/images/katya.jpg"),
		},
	}

	for _, dream := range dreams {
		database.Create(&dream)
		log.Printf("✅ Мечта: %s", dream.Title)
	}

	log.Println("🎉 Seed завершён!")
}

func strPtr(s string) *string {
	return &s
}

