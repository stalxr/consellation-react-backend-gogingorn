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

	// Создаем мечты (dreams) - расширенный набор для тестирования пагинации
	dreams := []models.Dream{
		{
			ID:              uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Title:           "Мечта Вани: увидеть море",
			Slug:            "vanya-more",
			ShortDescription: strPtr("Ваня мечтает увидеть море, но у него нет денег на поездку. Помогите осуществить мечту!"),
			FullDescription: strPtr("<p>Ване 10 лет, он живет в маленьком городе в Сибири и никогда не видел моря. Его мечта — хотя бы раз увидеть океан и побывать на пляже.</p><p>Сбор средств идет на поездку всей семьей в Сочи.</p>"),
			TargetAmount:    85000,
			CollectedAmount: 32500,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_vanya.jpg"),
			GalleryImages:   []string{"/uploads/images/vanya_1.jpg", "/uploads/images/vanya_2.jpg"},
		},
		{
			ID:              uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Title:           "Новый велосипед для Маши",
			Slug:            "masha-velosiped",
			ShortDescription: strPtr("Маша хочет кататься на велосипеде с друзьями и быть как все"),
			FullDescription: strPtr("<p>Маше 12 лет, она очень хочет велосипед, чтобы кататься с подругами в парке. Семья не может позволить себе покупку.</p>"),
			TargetAmount:    35000,
			CollectedAmount: 28000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_masha.jpg"),
			GalleryImages:   []string{"/uploads/images/masha_1.jpg"},
		},
		{
			ID:              uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			Title:           "Фотоаппарат для Сергея",
			Slug:            "sergey-fotoapparat",
			ShortDescription: strPtr("Сергей мечтает стать фотографом и запечатлеть красоту мира"),
			FullDescription: strPtr("<p>Сергей увлекается фотографией с детства. Его мечта — профессиональный фотоаппарат, чтобы развивать талант.</p>"),
			TargetAmount:    75000,
			CollectedAmount: 75000,
			Status:          "completed",
			CoverImage:      strPtr("/uploads/images/dream_sergey.jpg"),
			GalleryImages:   []string{"/uploads/images/sergey_1.jpg", "/uploads/images/sergey_2.jpg", "/uploads/images/sergey_3.jpg"},
		},
		{
			ID:              uuid.MustParse("44444444-4444-4444-4444-444444444444"),
			Title:           "Ноутбук для учебы Ани",
			Slug:            "anya-noutbuk",
			ShortDescription: strPtr("Ане нужен ноутбук для дистанционного обучения"),
			FullDescription: strPtr("<p>Аня отличница, но из-за пандемии нуждается в компьютере для онлайн-учебы. Семья не может купить.</p>"),
			TargetAmount:    45000,
			CollectedAmount: 12000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_anya.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("55555555-5555-5555-5555-555555555555"),
			Title:           "Лечение для котенка Мурзика",
			Slug:            "murzik-lechenie",
			ShortDescription: strPtr("Мурзику нужна операция, чтобы снова бегать и играть"),
			FullDescription: strPtr("<p>Мурзик попал в беду и сломал лапку. Операция дорогая, но без нее котенок не сможет нормально жить.</p>"),
			TargetAmount:    25000,
			CollectedAmount: 5000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_murzik.jpg"),
			GalleryImages:   []string{"/uploads/images/murzik_1.jpg"},
		},
		{
			ID:              uuid.MustParse("66666666-6666-6666-6666-666666666666"),
			Title:           "Гитара для Пети",
			Slug:            "petya-gitara",
			ShortDescription: strPtr("Петя мечтает играть на гитаре и писать песни"),
			FullDescription: strPtr("<p>Петя сочиняет стихи и мечтает научиться играть на гитаре, чтобы сочинять песни.</p>"),
			TargetAmount:    20000,
			CollectedAmount: 19500,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_petya.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("77777777-7777-7777-7777-777777777777"),
			Title:           "Поездка в зоопарк для класса",
			Slug:            "klass-zoopark",
			ShortDescription: strPtr("Дети из интерната хотят увидеть настоящих животных"),
			FullDescription: strPtr("<p>Дети из детского дома никогда не были в зоопарке. Помогите организовать поездку!</p>"),
			TargetAmount:    60000,
			CollectedAmount: 15000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_zoo.jpg"),
			GalleryImages:   []string{"/uploads/images/zoo_1.jpg"},
		},
		{
			ID:              uuid.MustParse("88888888-8888-8888-8888-888888888888"),
			Title:           "Курсы программирования для Коли",
			Slug:            "kolya-kursy",
			ShortDescription: strPtr("Коля хочет стать программистом и выйти из кризисной ситуации"),
			FullDescription: strPtr("<p>Коля самостоятельно учит Python, но нуждается в курсах с менторами для трудоустройства.</p>"),
			TargetAmount:    40000,
			CollectedAmount: 40000,
			Status:          "completed",
			CoverImage:      strPtr("/uploads/images/dream_kolya.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("99999999-9999-9999-9999-999999999999"),
			Title:           "Инвалидная коляска для бабушки",
			Slug:            "babushka-kolyaska",
			ShortDescription: strPtr("Новая коляска поможет бабушке выйти на прогулки"),
			FullDescription: strPtr("<p>Бабушка Лидия Ивановна прикована к постели. Новая коляска вернет ее к жизни.</p>"),
			TargetAmount:    55000,
			CollectedAmount: 8000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_babushka.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			Title:           "Спортивная форма для команды",
			Slug:            "komanda-forma",
			ShortDescription: strPtr("Детская футбольная команда нуждается в форме"),
			FullDescription: strPtr("<p>Ребята занимаются футболом, но у них нет денег на форму для турнира.</p>"),
			TargetAmount:    30000,
			CollectedAmount: 8000,
			Status:          "frozen",
			CoverImage:      strPtr("/uploads/images/dream_football.jpg"),
			GalleryImages:   []string{},
		},
		{
			ID:              uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			Title:           "Мебель в приют для животных",
			Slug:            "priyt-mebel",
			ShortDescription: strPtr("Приюту нужны клетки и лежанки для бездомных животных"),
			FullDescription: strPtr("<p>Маленький приют переполнен. Нужна мебель, чтобы спасти больше животных.</p>"),
			TargetAmount:    80000,
			CollectedAmount: 15000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_priyt.jpg"),
			GalleryImages:   []string{"/uploads/images/priyt_1.jpg"},
		},
		{
			ID:              uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			Title:           "Книги для школьной библиотеки",
			Slug:            "biblioteka-knigi",
			ShortDescription: strPtr("Сельской школе нужны новые учебники и художественная литература"),
			FullDescription: strPtr("<p>Школьная библиотека почти пуста. Дети жаждут читать, но книг нет.</p>"),
			TargetAmount:    25000,
			CollectedAmount: 7000,
			Status:          "active",
			CoverImage:      strPtr("/uploads/images/dream_books.jpg"),
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
