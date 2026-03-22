package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"charity-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB — глобальный объект подключения к базе данных (GORM).
var DB *gorm.DB

// Connect подключается к PostgreSQL,
// дожидается готовности контейнера через ретраи и выполняет миграции схемы.
func Connect() {
	var dsn string
	
	// Проверяем DATABASE_URL (Railway, Render и т.д.)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Railway может использовать DATABASE_PUBLIC_URL
		databaseURL = os.Getenv("DATABASE_PUBLIC_URL")
	}
	
	log.Printf("🔍 DATABASE_URL найден: %v", databaseURL != "")
	log.Printf("🔍 DATABASE_PUBLIC_URL найден: %v", os.Getenv("DATABASE_PUBLIC_URL") != "")
	
	if databaseURL != "" {
		dsn = databaseURL
		log.Println("🔌 Подключаемся к Postgres через DATABASE_URL...")
	} else {
		// Локальная разработка или Docker Compose
		host := "localhost"
		port := 5432
		if _, err := os.Stat("/.dockerenv"); err == nil {
			host = "db"
			port = 5432
		}
		dsn = fmt.Sprintf("host=%s user=postgres dbname=charity port=%d sslmode=disable TimeZone=UTC", host, port)
		log.Println("🔌 Подключаемся к Postgres (локально)...")
	}

	// Ретрай-луп: контейнер может подняться позже приложения.
	// 10 попыток по 2 секунды.
	var err error
	for i := 1; i <= 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("⏳ База ещё не готова (попытка %d/10): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("❌ Не удалось подключиться к базе данных: ", err)
	}

	log.Println("✅ Подключение к базе успешно!")

	// Миграции (структуры должны соответствовать OpenAPI).
	log.Println("📦 Запускаем миграции...")
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Dream{},
		&models.Donation{},
		&models.Report{},
		&models.News{},
	); err != nil {
		log.Fatal("❌ Ошибка миграций: ", err)
	}

	fmt.Println("✅ Миграции выполнены")
}
