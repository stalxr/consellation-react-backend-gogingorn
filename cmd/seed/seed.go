package main

import (
	"fmt"
	"log"
	"time"

	"charity-backend/internal/models"
	"charity-backend/pkg/db"
	"charity-backend/pkg/utils"

	"github.com/google/uuid"
)

// SeedData - заполнение базы тестовыми данными
func SeedData() {
	database := db.DB
	if database == nil {
		log.Fatal("❌ База данных не подключена")
	}

	log.Println("🌱 Начинаем заполнение тестовыми данными...")

	// Создаем админа
	adminPassword, _ := utils.HashPassword("admin123")
	admin := &models.User{
		ID:           uuid.New(),
		Email:        "admin@nastenka.ru",
		PasswordHash: adminPassword,
		FullName:     "Администратор",
		Role:         "admin",
		TotalDonated: 0,
	}

	// Проверяем существование админа
	var existingAdmin models.User
	if err := database.Where("email = ?", admin.Email).First(&existingAdmin).Error; err != nil {
		if err := database.Create(admin).Error; err != nil {
			log.Printf("❌ Ошибка создания админа: %v", err)
		} else {
			log.Println("✅ Админ создан: admin@nastenka.ru / admin123")
		}
	} else {
		log.Println("ℹ️ Админ уже существует")
	}

	// Создаем тестового пользователя
	userPassword, _ := utils.HashPassword("user123")
	user := &models.User{
		ID:           uuid.New(),
		Email:        "user@example.com",
		PasswordHash: userPassword,
		FullName:     "Тестовый Пользователь",
		Role:         "user",
		TotalDonated: 1500,
	}

	var existingUser models.User
	if err := database.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		if err := database.Create(user).Error; err != nil {
			log.Printf("❌ Ошибка создания пользователя: %v", err)
		} else {
			log.Println("✅ Пользователь создан: user@example.com / user123")
		}
	} else {
		log.Println("ℹ️ Пользователь уже существует")
	}

	// Создаем мечты (dreams) - дети с онкологическими заболеваниями
	dreams := []models.Dream{
		{
			ID:              uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Title:           "Лечение лейкемии для Маши, 7 лет",
			Slug:            "lechenie-leykemii-mashi-7-let",
			ShortDescription: strPtr("Маше нужен курс химиотерапии и пересадка костного мозга"),
			FullDescription: strPtr("<p>У Маши диагностирована острая лимфобластная лейкемия. Необходим курс химиотерапии (6 месяцев) с последующей пересадкой костного мозга. Стоимость лечения в НМИЦ онкологии — 2,800,000 рублей.</p><p>Семья уже израсходовала все сбережения на первый курс химиотерапии.</p>"),
			TargetAmount:    2800000,
			CollectedAmount: 950000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/masha_oncology.jpg"),
			GalleryImages:   []string{"/uploads/images/masha_1.jpg", "/uploads/images/masha_2.jpg"},
		},
		{
			ID:              uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Title:           "Операция по удалению опухоли мозга для Вани, 5 лет",
			Slug:            "operaciya-opuhol-mozga-vanya-5-let",
			ShortDescription: strPtr("Сложная нейрохирургическая операция в Израиле"),
			FullDescription: strPtr("<p>У Вани обнаружена злокачественная опухоль головного мозга (медуллобластома). Операцию можно провести только в специализированном центре в Израиле. Стоимость: 3,500,000 рублей.</p><p>Без операции ребенок не выживет.</p>"),
			TargetAmount:    3500000,
			CollectedAmount: 2100000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/vanya_brain.jpg"),
			GalleryImages:   []string{"/uploads/images/vanya_1.jpg"},
		},
		{
			ID:              uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			Title:           "Протонная терапия для Кати, 12 лет (рак глаза)",
			Slug:            "protonnaya-terapiya-katya-12-let",
			ShortDescription: strPtr("Современное лечение ретинобластомы в Германии"),
			FullDescription: strPtr("<p>У Кати редкая форма рака сетчатки глаза — ретинобластома. Единственный шанс сохранить зрение — протонная терапия в клинике Гейдельберга, Германия. Стоимость курса: 1,800,000 рублей.</p>"),
			TargetAmount:    1800000,
			CollectedAmount: 1800000,
			Status:          "completed",
			CoverImage:      strPtr("/uploads/images/katya_eye.jpg"),
			GalleryImages:   []string{"/uploads/images/katya_1.jpg", "/uploads/images/katya_2.jpg"},
		},
		{
			ID:              uuid.MustParse("44444444-4444-4444-4444-444444444444"),
			Title:           "Таргетная терапия для Димы, 9 лет (саркома)",
			Slug:            "target-terapiya-dima-9-let",
			ShortDescription: strPtr("Иммунотерапия для борьбы с остеосаркомой"),
			FullDescription: strPtr("<p>У Димы диагностирована остеосаркома (злокачественная опухоль кости). Стандартная химиотерапия неэффективна, необходима таргетная иммунотерапия препаратом Динутуксимаб бета. Курс: 1,200,000 рублей.</p>"),
			TargetAmount:    1200000,
			CollectedAmount: 450000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dima_sarcoma.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("55555555-5555-5555-5555-555555555555"),
			Title:           "Реабилитация после лучевой терапии для Насти, 6 лет",
			Slug:            "reabilitaciya-nastya-6-let",
			ShortDescription: strPtr("Восстановление после лечения нейробластомы"),
			FullDescription: strPtr("<p>Настя успешно прошла курс лучевой терапии от нейробластомы. Теперь необходима длительная реабилитация в центре 'Алиса' для восстановления иммунитета и нервной системы. Стоимость: 680,000 рублей.</p>"),
			TargetAmount:    680000,
			CollectedAmount: 320000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/nastya_rehab.jpg"),
			GalleryImages:   []string{"/uploads/images/nastya_1.jpg"},
		},
		{
			ID:              uuid.MustParse("66666666-6666-6666-6666-666666666666"),
			Title:           "Гамма-нож для Лёши, 14 лет (глиома ствола мозга)",
			Slug:            "gamma-nozh-lesha-14-let",
			ShortDescription: strPtr("Бескровная операция на стволе головного мозга"),
			FullDescription: strPtr("<p>У Лёши неоперабельная глиома ствола мозга. Единственный метод — радиохирургия 'Гамма-нож' в Москве. Процедура позволит остановить рост опухоли. Стоимость: 850,000 рублей.</p>"),
			TargetAmount:    850000,
			CollectedAmount: 0,
			Status:          "pending",
			CoverImage:      strPtr("/uploads/images/lesha_gamma.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("77777777-7777-7777-7777-777777777777"),
			Title:           "CAR-T терапия для Полины, 11 лет (рецидив лейкемии)",
			Slug:            "car-t-terapiya-polina-11-let",
			ShortDescription: strPtr("Генная терапия для борьбы с рецидивом"),
			FullDescription: strPtr("<p>Полина перенесла лейкемию, но болезнь вернулась. Обычная химиотерапия бессильна. CAR-T терапия — последний шанс. Стоимость инновационного лечения: 5,200,000 рублей.</p>"),
			TargetAmount:    5200000,
			CollectedAmount: 1800000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/polina_cart.jpg"),
			GalleryImages:   []string{"/uploads/images/polina_1.jpg"},
		},
		{
			ID:              uuid.MustParse("88888888-8888-8888-8888-888888888888"),
			Title:           "Трансплантация печени для Саши, 3 года (гепатобластома)",
			Slug:            "transplantaciya-pecheni-sasha-3-goda",
			ShortDescription: strPtr("Срочная пересадка печени в Турции"),
			FullDescription: strPtr("<p>У Саши раковая опухоль печени — гепатобластома. Опухоль занимает 70% органа, химиотерапия не помогает. Необходима срочная трансплантация в клинике Стамбула. Стоимость: 4,800,000 рублей.</p>"),
			TargetAmount:    4800000,
			CollectedAmount: 4800000,
			Status:          "completed",
			CoverImage:      strPtr("/uploads/images/sasha_liver.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("99999999-9999-9999-9999-999999999999"),
			Title:           "Диагностика ПЭТ-КТ для Миши, 8 лет (подозрение на лимфому)",
			Slug:            "pet-kt-misha-8-let",
			ShortDescription: strPtr("Точная диагностика для назначения лечения"),
			FullDescription: strPtr("<p>У Миши подозрение на лимфому Ходжкина. Для точной стадии и назначения правильного лечения необходим ПЭТ-КТ сканер. Исследование: 180,000 рублей.</p>"),
			TargetAmount:    180000,
			CollectedAmount: 45000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/misha_pet.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			Title:           "Поддерживающая терапия для Вари, 16 лет (рецидив)",
			Slug:            "poddervzhivayushchaya-terapiya-vara-16-let",
			ShortDescription: strPtr("Паллиативная помощь и обезболивание"),
			FullDescription: strPtr("<p>Варя борется с раком уже 4 года. Сейчас необходима качественная паллиативная помощь: обезболивание, питание, уход. Ежемесячные расходы: 120,000 рублей.</p>"),
			TargetAmount:    120000,
			CollectedAmount: 35000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/varia_palliative.jpg"),
			GalleryImages:   []string{},
		},
	}

	for _, dream := range dreams {
		var existingDream models.Dream
		if err := database.Where("slug = ?", dream.Slug).First(&existingDream).Error; err != nil {
			if err := database.Create(&dream).Error; err != nil {
				log.Printf("❌ Ошибка создания мечты %s: %v", dream.Title, err)
			} else {
				log.Printf("✅ Мечта создана: %s", dream.Title)
			}
		} else {
			log.Printf("ℹ️ Мечта %s уже существует", dream.Title)
		}
	}

	// Создаем пожертвования для тестирования истории и списка доноров
	donations := []models.Donation{
		// Пожертвования для мечты Вани
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserID:      &user.ID,
			Amount:      5000,
			Status:      "completed",
			IsAnonymous: false,
			Comment:     strPtr("Удачи, Ваня!"),
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserID:      nil,
			Amount:      3000,
			Status:      "completed",
			IsAnonymous: true,
			Email:       strPtr("anon@example.com"),
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserID:      nil,
			Amount:      10000,
			Status:      "completed",
			IsAnonymous: false,
			Email:       strPtr("maria@example.com"),
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserID:      &user.ID,
			Amount:      1500,
			Status:      "pending",
			IsAnonymous: false,
		},
		// Пожертвования для мечты Маши
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			UserID:      &user.ID,
			Amount:      2000,
			Status:      "completed",
			IsAnonymous: false,
			Comment:     strPtr("Катайся с ветерком!"),
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			UserID:      nil,
			Amount:      5000,
			Status:      "completed",
			IsAnonymous: true,
		},
		// Пожертвования для завершенной мечты Сергея
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			UserID:      &user.ID,
			Amount:      15000,
			Status:      "completed",
			IsAnonymous: false,
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			UserID:      nil,
			Amount:      20000,
			Status:      "completed",
			IsAnonymous: false,
			Email:       strPtr("photolover@example.com"),
			Comment:     strPtr("Жду твои фотографии!"),
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			UserID:      nil,
			Amount:      10000,
			Status:      "completed",
			IsAnonymous: true,
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			UserID:      nil,
			Amount:      5000,
			Status:      "completed",
			IsAnonymous: false,
			Email:       strPtr("helper@example.com"),
		},
		// Пожертвования для мечты Ани
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("44444444-4444-4444-4444-444444444444"),
			UserID:      nil,
			Amount:      3000,
			Status:      "completed",
			IsAnonymous: false,
			Email:       strPtr("teacher@school.ru"),
			Comment:     strPtr("Учись на отлично!"),
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("44444444-4444-4444-4444-444444444444"),
			UserID:      &user.ID,
			Amount:      5000,
			Status:      "completed",
			IsAnonymous: false,
		},
		// Небольшие пожертвования для других мечт
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("55555555-5555-5555-5555-555555555555"),
			UserID:      nil,
			Amount:      1000,
			Status:      "completed",
			IsAnonymous: true,
		},
		{
			ID:          uuid.New(),
			DreamID:     uuid.MustParse("66666666-6666-6666-6666-666666666666"),
			UserID:      nil,
			Amount:      500,
			Status:      "completed",
			IsAnonymous: false,
			Email:       strPtr("musicfan@example.com"),
		},
	}

	// Создаем пожертвования
	donationCount := 0
	for _, donation := range donations {
		var existingDonation models.Donation
		if err := database.Where("id = ?", donation.ID).First(&existingDonation).Error; err != nil {
			if err := database.Create(&donation).Error; err != nil {
				log.Printf("❌ Ошибка создания пожертвования: %v", err)
			} else {
				donationCount++
			}
		}
	}
	if donationCount > 0 {
		log.Printf("✅ Создано пожертвований: %d", donationCount)
	}

	// Создаем новости с датами публикации
	now := time.Now()
	news := []models.News{
		{
			ID:          uuid.New(),
			Title:       "Фонд помог 100 детям в этом году",
			PreviewText: strPtr("Благодаря вашей поддержке мы смогли помочь сотне детей осуществить их мечты..."),
			Content:     strPtr("<p>За этот год благодаря вашей поддержке мы смогли помочь 100 детям из разных уголков России. Каждая мечта, каждая улыбка ребенка — это благодаря вам.</p><p>Подробная статистика и истории успеха доступны в нашем годовом отчете.</p>"),
			ImageURL:    strPtr("/uploads/images/news1.jpg"),
			PublishedAt: &now,
		},
		{
			ID:          uuid.New(),
			Title:       "Новогодняя акция 'Созвездие мечты'",
			PreviewText: strPtr("Приглашаем всех принять участие в нашей новогодней акции помощи детям..."),
			Content:     strPtr("<p>В период новогодних праздников мы запускаем специальную акцию. Сделайте подарок ребенку, помогите осуществить мечту.</p><p>Каждое пожертвование в декабре будет удвоено нашим партнером.</p>"),
			ImageURL:    strPtr("/uploads/images/news2.jpg"),
			PublishedAt: &[]time.Time{now.AddDate(0, -1, 0)}[0],
		},
		{
			ID:          uuid.New(),
			Title:       "Открыт новый центр помощи в Москве",
			PreviewText: strPtr("Мы открыли новый офис, где можно лично познакомиться с нашей работой..."),
			Content:     strPtr("<p>Теперь каждый желающий может посетить наш центр, познакомиться с командой и узнать больше о наших проектах.</p>"),
			ImageURL:    strPtr("/uploads/images/news3.jpg"),
			PublishedAt: &[]time.Time{now.AddDate(0, -2, 0)}[0],
		},
	}

	for _, n := range news {
		var existingNews models.News
		if err := database.Where("title = ?", n.Title).First(&existingNews).Error; err != nil {
			if err := database.Create(&n).Error; err != nil {
				log.Printf("❌ Ошибка создания новости %s: %v", n.Title, err)
			} else {
				log.Printf("✅ Новость создана: %s", n.Title)
			}
		} else {
			log.Printf("ℹ️ Новость %s уже существует", n.Title)
		}
	}

	// Создаем финансовые отчеты
	reports := []models.Report{
		{
			ID:      uuid.New(),
			Year:    now.Year(),
			Month:   int(now.Month()),
			Title:   fmt.Sprintf("Финансовый отчет за %s %d", now.Month().String(), now.Year()),
			FileURL: "/uploads/reports/report_current.pdf",
		},
		{
			ID:      uuid.New(),
			Year:    now.Year(),
			Month:   int(now.Month()) - 1,
			Title:   fmt.Sprintf("Финансовый отчет за %s %d", now.AddDate(0, -1, 0).Month().String(), now.Year()),
			FileURL: "/uploads/reports/report_last_month.pdf",
		},
		{
			ID:      uuid.New(),
			Year:    now.Year() - 1,
			Month:   12,
			Title:   fmt.Sprintf("Годовой отчет за %d год", now.Year()-1),
			FileURL: "/uploads/reports/report_annual.pdf",
		},
		{
			ID:      uuid.New(),
			Year:    now.Year(),
			Month:   1,
			Title:   fmt.Sprintf("Финансовый отчет за Январь %d", now.Year()),
			FileURL: "/uploads/reports/report_jan.pdf",
		},
	}

	reportCount := 0
	for _, r := range reports {
		var existingReport models.Report
		if err := database.Where("id = ?", r.ID).First(&existingReport).Error; err != nil {
			if err := database.Create(&r).Error; err != nil {
				log.Printf("❌ Ошибка создания отчета: %v", err)
			} else {
				reportCount++
			}
		}
	}
	if reportCount > 0 {
		log.Printf("✅ Создано отчетов: %d", reportCount)
	}

	log.Println("🎉 Заполнение данными завершено!")
	log.Println("")
	log.Println("🔑 Тестовые учетные записи:")
	log.Println("   Админ: admin@nastenka.ru / admin123")
	log.Println("   Юзер:  user@example.com / user123")
	log.Println("")
	log.Println("📚 API доступно по адресу: http://localhost:8080/api/v1")
	log.Println("📖 Документация Swagger: http://localhost:8080/swagger/index.html")
}

func strPtr(s string) *string {
	return &s
}

// main для запуска seed данных
func main() {
	// Подключаемся к БД
	db.Connect()
	
	// Заполняем данными
	SeedData()
}
