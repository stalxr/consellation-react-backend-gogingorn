package db

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	
	// Debug: покажем все переменные окружения
	log.Println("🔍 Все переменные окружения:")
	for _, env := range os.Environ() {
		if strings.Contains(env, "PG") || strings.Contains(env, "DATABASE") || strings.Contains(env, "POSTGRES") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				log.Printf("  %s=%s", parts[0], parts[1])
			}
		}
	}
	
	// Проверяем DATABASE_URL (Railway, Render и т.д.)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_PUBLIC_URL")
	}
	
	// Railway предоставляет отдельные переменные PGHOST, PGPORT и т.д.
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")
	pgUser := os.Getenv("PGUSER")
	pgPassword := os.Getenv("PGPASSWORD")
	pgDatabase := os.Getenv("PGDATABASE")
	
	// Альтернативные имена переменных
	if pgHost == "" {
		pgHost = os.Getenv("POSTGRES_HOST")
	}
	if pgPort == "" {
		pgPort = os.Getenv("POSTGRES_PORT")
	}
	if pgUser == "" {
		pgUser = os.Getenv("POSTGRES_USER")
	}
	if pgPassword == "" {
		pgPassword = os.Getenv("POSTGRES_PASSWORD")
	}
	if pgDatabase == "" {
		pgDatabase = os.Getenv("POSTGRES_DB")
	}
	
	log.Printf("🔍 DATABASE_URL найден: %v", databaseURL != "")
	log.Printf("🔍 PGHOST: '%s', PGPORT: '%s', PGUSER: '%s', PGDATABASE: '%s'", pgHost, pgPort, pgUser, pgDatabase)
	
	if databaseURL != "" {
		dsn = databaseURL
		log.Println("🔌 Подключаемся через DATABASE_URL...")
	} else if pgHost != "" && pgPort != "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			pgHost, pgUser, pgPassword, pgDatabase, pgPort)
		log.Println("🔌 Подключаемся через Railway PG* переменные...")
	} else {
		host := "localhost"
		port := 5432
		if _, err := os.Stat("/.dockerenv"); err == nil {
			host = "db"
		}
		dsn = fmt.Sprintf("host=%s user=postgres dbname=charity port=%d sslmode=disable TimeZone=UTC", host, port)
		log.Println("🔌 Подключаемся локально...")
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
