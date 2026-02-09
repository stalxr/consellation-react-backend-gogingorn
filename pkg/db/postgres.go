package db

import (
	"fmt"
	"log"
	"os"

	"charity-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB подключается к базе данных PostgreSQL и выполняет автоматические миграции
// Читает настройки подключения из переменных окружения (.env файл)
// Возвращает объект *gorm.DB для работы с базой данных или ошибку при неудаче
func ConnectDB() (*gorm.DB, error) {
	// Читаем настройки из переменных окружения
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "charity_db"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	// Формируем DSN строку
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	// Подключаемся к БД
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	// Включаем расширение pgcrypto, чтобы работал gen_random_uuid()
	if err := database.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		return nil, fmt.Errorf("не удалось включить расширение pgcrypto: %w", err)
	}

	// Выполняем автоматическую миграцию всех моделей
	err = database.AutoMigrate(
		&models.User{},
		&models.Dream{},
		&models.Donation{},
		&models.News{},
		&models.Report{},
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка миграции: %w", err)
	}

	log.Println("✅ База данных успешно подключена и миграции выполнены")
	return database, nil
}

