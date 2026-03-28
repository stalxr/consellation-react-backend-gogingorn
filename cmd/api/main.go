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
	// force=true — очистить старые данные и создать новые 25 мечт
	SeedData(database, true)

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
// При force=true очищает старые данные и создаёт новые
func SeedData(database *gorm.DB, force bool) {
	log.Println("🌱 Начинаем заполнение тестовыми данными...")

	if force {
		// Очищаем старые данные
		log.Println("🗑️ Очищаем старые данные...")
		database.Exec("DELETE FROM donations")
		database.Exec("DELETE FROM dreams")
		database.Exec("DELETE FROM users WHERE email IN ('admin@nastenka.ru', 'user@example.com')")
	}

	// Проверяем есть ли данные
	var count int64
	database.Model(&models.Dream{}).Count(&count)
	if count > 0 && !force {
		log.Println("ℹ️ Данные уже существуют, пропускаем seed (используйте force=true для пересоздания)")
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

	// 25 детей с мечтами
	dreams := []models.Dream{
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000001"),
			Title:            "Тимоша, 4 года — полететь на самолете на море",
			Slug:             "timosha-4-goda-poletet-na-samoletye-na-more",
			ShortDescription: strPtr("Мечтаю, что когда я буду здоров, обязательно полечу на настоящем самолете на море"),
			FullDescription:  strPtr("<p>Мечтаю, что когда я буду здоров, обязательно полечу на настоящем самолете на море. Со мной будут рядом мама, папа, сестренки и братик. История Тимоши — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     2500000,
			CollectedAmount:  850000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/timosha.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000002"),
			Title:            "Давид, 6 лет — стать пожарным и спасти людей",
			Slug:             "david-6-let-stat-pozharnym",
			ShortDescription: strPtr("Я мечтаю поскорее выздороветь и стать пожарным"),
			FullDescription:  strPtr("<p>«Я мечтаю поскорее выздороветь и стать пожарным, чтобы спасти как можно больше людей!» История Давида — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     1800000,
			CollectedAmount:  620000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/david.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000003"),
			Title:            "Вася, 16 лет — стать регулировщиком",
			Slug:             "vasya-16-let-stat-regulirovshikom",
			ShortDescription: strPtr("Я могу часами наблюдать за движением на перекрёстке"),
			FullDescription:  strPtr("<p>«Я могу часами наблюдать за движением на перекрёстке. Жду того дня, когда смогу успешно сдать экзамен и наконец-то занять место регулировщика!» История Васи — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     1200000,
			CollectedAmount:  450000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/vasya.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000004"),
			Title:            "Катя, 4 года — увидеть живого ёжика",
			Slug:             "katya-4-goda-uvidet-yozhika",
			ShortDescription: strPtr("Хочу увидеть живого ёжика!"),
			FullDescription:  strPtr("<p>«Хочу увидеть живого ёжика! Когда я смогу ходить самостоятельно, мы с мамой пойдем на прогулку в лес и обязательно его увидим.» История Кати — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     1500000,
			CollectedAmount:  980000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/katya_yozhik.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000005"),
			Title:            "Дима, 9 лет — стать фотографом",
			Slug:             "dima-9-let-stat-fotografom",
			ShortDescription: strPtr("Мечтаю выздороветь и стать фотографом"),
			FullDescription:  strPtr("<p>«Мечтаю выздороветь и стать фотографом. Чтобы показать людям всю красоту мира, которую я пока могу видеть только из окна больничной палаты.» История Димы — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     2100000,
			CollectedAmount:  1500000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/dima_foto.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000006"),
			Title:            "Лёня, 11 лет — стать футболистом",
			Slug:             "lenya-11-let-stat-futbolistom",
			ShortDescription: strPtr("Моя мечта — стать известным футболистом"),
			FullDescription:  strPtr("<p>«Моя мечта — стать известным футболистом. Когда я выздоровею, то обязательно стану лучшим игроком!» История Леонида — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     3200000,
			CollectedAmount:  2100000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/lenya.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000007"),
			Title:            "Эмилия, 7 лет — увидеть настоящие горы",
			Slug:             "emiliya-7-let-uvidet-gory",
			ShortDescription: strPtr("Я мечтаю увидеть настоящие горы"),
			FullDescription:  strPtr("<p>«Я мечтаю увидеть настоящие горы. Когда я буду здорова мы с мамой пойдем в поход по горным тропкам.» История Эмилии — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     2800000,
			CollectedAmount:  1200000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/emiliya.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000008"),
			Title:            "Дамир, 10 лет — увидеть улыбку мамы",
			Slug:             "damir-10-let-ulibku-mamy",
			ShortDescription: strPtr("Хочу увидеть счастливую улыбку мамы"),
			FullDescription:  strPtr("<p>«Я очень хочу увидеть счастливую улыбку мамы, когда она узнает, что я здоров. Хочу обнимать её долго-долго и больше никогда не разлучаться.» История Дамира — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     1900000,
			CollectedAmount:  750000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/damir.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000009"),
			Title:            "Ксюша, 6 лет — есть конфеты и мороженое",
			Slug:             "ksyusha-6-let-konfety",
			ShortDescription: strPtr("Смогу есть всё, что захочу"),
			FullDescription:  strPtr("<p>«Когда я стану здоровой. Смогу есть всё, что захочу. Много-много конфет, а ещё мороженое и даже чипсы!» История Ксюши — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     1600000,
			CollectedAmount:  1600000,
			Status:           "completed",
			CoverImage:       strPtr("/uploads/images/ksyusha.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000010"),
			Title:            "Канышай, 15 лет — стать кондитером",
			Slug:             "kanyshay-15-let-konditer",
			ShortDescription: strPtr("Мечтаю выздороветь и стать кондитером"),
			FullDescription:  strPtr("<p>«Я мечтаю выздороветь и стать кондитером. Буду печь маме и друзьям самые вкусные и красивые торты.» История Канышай — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     2300000,
			CollectedAmount:  890000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/kanyshay.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000011"),
			Title:            "София, 14 лет — пойти на выпускной",
			Slug:             "sofiya-14-let-vypusknoy",
			ShortDescription: strPtr("Снова научиться ходить без поддержки"),
			FullDescription:  strPtr("<p>«Мечтаю снова научиться ходить без поддержки и пойти на школьный выпускной в красивых туфлях на высоком каблуке.» История Софии — всего одна звёздочка - мечта из множества, что сияет в этом созвездии.</p>"),
			TargetAmount:     3400000,
			CollectedAmount:  1200000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/sofiya.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000012"),
			Title:            "Гульнара, 10 лет — завести щеночка",
			Slug:             "gulnara-10-let-shchenok",
			ShortDescription: strPtr("Мечтаю о пушистом щеночке"),
			FullDescription:  strPtr("<p>«Мечтаю о пушистом щеночке. Когда я смогу ходить сама, мы будем с ним много гулять и играть!»</p>"),
			TargetAmount:     1400000,
			CollectedAmount:  0,
			Status:           "pending",
			CoverImage:       strPtr("/uploads/images/gulnara.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000013"),
			Title:            "Никита, 6 лет — завести котёнка",
			Slug:             "nikita-6-let-kotyenok",
			ShortDescription: strPtr("Я так хочу побыстрее выздороветь"),
			FullDescription:  strPtr("<p>«Я мечтаю завести маленького пушистого котёнка. Но мы с мамой почти всё время лежим в больнице, а кошечек сюда не пускают. Я так хочу побыстрее выздороветь!»</p>"),
			TargetAmount:     1800000,
			CollectedAmount:  650000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/nikita.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000014"),
			Title:            "Ваня, 12 лет — стать хоккеистом",
			Slug:             "vanya-12-let-hokkeist",
			ShortDescription: strPtr("Мечтаю стать известным хоккеистом"),
			FullDescription:  strPtr("<p>«Я мечтаю стать известным хоккеистом. Когда я выздоровею, то буду тренироваться каждый день!»</p>"),
			TargetAmount:     4100000,
			CollectedAmount:  1800000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/vanya_hockey.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000015"),
			Title:            "Матвей, 7 лет — стать блогером",
			Slug:             "matvey-7-let-bloger",
			ShortDescription: strPtr("Хочу стать известным блогером"),
			FullDescription:  strPtr("<p>«Я хочу стать известным блогером, чтобы снимать добрые видео и помогать другим детям, как сейчас помогают мне.»</p>"),
			TargetAmount:     1300000,
			CollectedAmount:  920000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/matvey.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000016"),
			Title:            "Дима, 7 лет — стать машинистом паровоза",
			Slug:             "dima-7-let-mashinist",
			ShortDescription: strPtr("Моя мечта — стать машинистом паровоза"),
			FullDescription:  strPtr("<p>«Моя мечта — стать машинистом паровоза. Но это очень ответственная работа, поэтому мне нужно быть внимательным и здоровым.»</p>"),
			TargetAmount:     1700000,
			CollectedAmount:  0,
			Status:           "pending",
			CoverImage:       strPtr("/uploads/images/dima_parovoz.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000017"),
			Title:            "Ярослав, 12 лет — ухаживать за питомцами",
			Slug:             "yaroslav-12-let-pitomci",
			ShortDescription: strPtr("У меня есть два лучших друга"),
			FullDescription:  strPtr("<p>«У меня есть целых два лучших друга — собачка Соня и попугай Веня! Хочу набраться побольше сил, чтобы хорошо за ними ухаживать.»</p>"),
			TargetAmount:     2200000,
			CollectedAmount:  1100000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/yaroslav.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000018"),
			Title:            "Алисия, 7 лет — увидеть черепашек",
			Slug:             "alisiya-7-let-cherepashki",
			ShortDescription: strPtr("Я так хочу увидеть забавных черепашек"),
			FullDescription:  strPtr("<p>«Когда я буду здорова, мы с мамой поедем на море. Я так хочу увидеть забавных, маленьких черепашек!»</p>"),
			TargetAmount:     2600000,
			CollectedAmount:  3400000,
			Status:           "completed",
			CoverImage:       strPtr("/uploads/images/alisiya.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000019"),
			Title:            "Лиза, 9 лет — научиться кататься на коньках",
			Slug:             "liza-9-let-konki",
			ShortDescription: strPtr("Зимой на катке так весело"),
			FullDescription:  strPtr("<p>«Зимой на катке так весело. Когда я буду здоровой мы обязательно туда поедем, и я научусь кататься на коньках.»</p>"),
			TargetAmount:     1100000,
			CollectedAmount:  480000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/liza.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000020"),
			Title:            "Мехродж, 17 лет — путешествие на корабле",
			Slug:             "mehrodzh-17-let-korabl",
			ShortDescription: strPtr("Море, ветер и парус — это моя мечта"),
			FullDescription:  strPtr("<p>«Море, ветер и парус — это и есть моя мечта! Когда я стану здоровым хочу отправиться в путешествие на корабле и увидеть много красивых мест.»</p>"),
			TargetAmount:     3800000,
			CollectedAmount:  1500000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/mehrodzh.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000021"),
			Title:            "Ксюша, 13 лет — стать врачом",
			Slug:             "ksyusha-13-let-vrach",
			ShortDescription: strPtr("Моя мечта — стать врачом и спасти жизни"),
			FullDescription:  strPtr("<p>«Моя мечта — стать врачом и спасти множество жизней. Но сначала мне нужно выздороветь и набраться сил, чтобы помогать другим.»</p>"),
			TargetAmount:     2900000,
			CollectedAmount:  2100000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/ksyusha_vrach.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000022"),
			Title:            "Рома, 7 лет — полетать на вертолете",
			Slug:             "roma-7-let-vertolyot",
			ShortDescription: strPtr("Я мечтаю полетать на вертолете"),
			FullDescription:  strPtr("<p>«Я мечтаю полетать на вертолете! А ещё когда я стану здоровым и вырасту, обязательно научусь им управлять.»</p>"),
			TargetAmount:     3500000,
			CollectedAmount:  2800000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/roma.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000023"),
			Title:            "Мирослава, 7 лет — обнять свою куклу",
			Slug:             "miroslava-7-let-kukla",
			ShortDescription: strPtr("Я так скучаю по своей кукле"),
			FullDescription:  strPtr("<p>«Я так скучаю по своей кукле! Моя мечта — побыстрее выздороветь и поехать домой, чтобы её обнять.»</p>"),
			TargetAmount:     900000,
			CollectedAmount:  900000,
			Status:           "completed",
			CoverImage:       strPtr("/uploads/images/miroslava.jpg"),
		},
		{
			ID:               uuid.MustParse("a0000000-0000-0000-0000-000000000024"),
			Title:            "Ульяна, 5 лет — велосипед",
			Slug:             "ulyana-5-let-velosiped",
			ShortDescription: strPtr("Хочу велосипед! Большой и красивый"),
			FullDescription:  strPtr("<p>«Хочу велосипед! Большой и красивый. Но я ещё не умею на нём кататься. Сделаю это, когда разрешат врачи.»</p>"),
			TargetAmount:     600000,
			CollectedAmount:  230000,
			Status:           "active",
			CoverImage:       strPtr("/uploads/images/ulyana.jpg"),
		},
	}

	for _, dream := range dreams {
		database.Create(&dream)
		log.Printf("✅ Мечта: %s", dream.Title)
	}

	log.Println("🎉 Seed завершён! Создано", len(dreams), "мечт")
}

func strPtr(s string) *string {
	return &s
}

