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

// Connect подключается к PostgreSQL (trust mode, без пароля),
// дожидается готовности контейнера через ретраи и выполняет миграции схемы.
func Connect() {
	// ВАЖНО: в Docker Compose мы используем trust-auth (POSTGRES_HOST_AUTH_METHOD=trust),
	// поэтому пароль в DSN НЕ указываем. Это сделано специально, чтобы на Windows
	// уйти от SASL/SCRAM циклов при локальной разработке.
	//
	// DSN строго по требованиям (для локального запуска):
	// host=localhost user=postgres dbname=charity port=5432 sslmode=disable TimeZone=UTC
	//
	// НЮАНС: если мы запускаем API в Docker, то "localhost" указывает на контейнер API,
	// а Postgres живёт в контейнере "db". Поэтому внутри Docker меняем host на "db".
	host := "localhost"
	port := 5432
	if _, err := os.Stat("/.dockerenv"); err == nil {
		host = "db"
		port = 5432
	}

	dsn := fmt.Sprintf("host=%s user=postgres dbname=charity port=%d sslmode=disable TimeZone=UTC", host, port)

	log.Println("🔌 Подключаемся к Postgres (trust auth, без пароля)...")

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
