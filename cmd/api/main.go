package main

import (
	"log"
	"os"
	"path/filepath"

	"charity-backend/internal/handlers"
	"charity-backend/internal/middleware"
	"charity-backend/internal/repository"
	"charity-backend/internal/services"
	"charity-backend/pkg/db"
	"charity-backend/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// Настраиваем статическую раздачу файлов из папки uploads
	// Это позволяет открывать загруженные изображения и PDF через браузер
	router.Static("/uploads", "./uploads")

	// API v1 группа
	api := router.Group("/api/v1")
	{
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

