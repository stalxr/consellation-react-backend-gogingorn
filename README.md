# Благотворительная платформа - Бэкенд API

## Структура проекта

```
consellation-react-backend-gogingorn/
├── cmd/
│   ├── api/
│   │   └── main.go              # Точка входа API
│   └── seed/
│       └── seed.go              # Скрипт мок-данных
├── internal/
│   ├── handlers/                # HTTP обработчики (Gin controllers)
│   │   ├── auth_handler.go      # Авторизация
│   │   ├── dream_handler.go     # Мечты (CRUD)
│   │   ├── donation_handler.go  # Пожертвования
│   │   ├── news_handler.go      # Новости
│   │   ├── payment_handler.go   # Платежи
│   │   ├── report_handler.go    # Отчеты
│   │   └── upload_handler.go    # Загрузка файлов
│   ├── services/                # Бизнес-логика
│   ├── repository/              # Работа с БД (GORM)
│   ├── middleware/              # JWT, CORS, Auth
│   └── models/                  # Структуры данных
├── pkg/
│   └── db/
│       └── postgres.go          # Подключение к PostgreSQL
├── openapi.yaml                 # Спецификация API (OpenAPI 3.0)
├── API_DOCUMENTATION.md         # Подробная документация
├── FRONTEND.md                  # Гайд для фронтендера ⭐
├── Dream_API_Postman_Collection.json
├── docker-compose.yaml          # PostgreSQL локально
├── Dockerfile                   # Сборка контейнера
├── go.mod                       # Go зависимости
└── .env                         # Конфигурация (пример)
```

## Быстрый старт для фронтендера

```bash
# 1. Клонируй
git clone https://github.com/stalxr/consellation-react-backend-gogingorn.git
cd consellation-react-backend-gogingorn

# 2. Запусти PostgreSQL
docker-compose up -d

# 3. Запусти API
go run cmd/api/main.go

# 4. (Опционально) Заполни мок-данными
go run cmd/seed/seed.go
```

API будет на `http://localhost:8080/api/v1`

## Документация для фронтенда

- **`FRONTEND.md`** - Гайд для фронтендера (как запустить, тестовые данные, примеры запросов)
- **`openapi.yaml`** - Спецификация API для Swagger/Postman
- **`API_DOCUMENTATION.md`** - Полная документация всех эндпоинтов
- **`Dream_API_Postman_Collection.json`** - Готовая коллекция Postman

## Тестовые аккаунты (после seed)

| Роль | Email | Пароль |
|------|-------|--------|
| Админ | admin@nastenka.ru | admin123 |
| Юзер | user@example.com | user123 |

## Основные эндпоинты

```
POST /api/v1/auth/register       # Регистрация
POST /api/v1/auth/login          # Вход
GET  /api/v1/auth/me             # Профиль (нужен Bearer токен)

GET  /api/v1/dreams              # Список мечт
GET  /api/v1/dreams/{id}         # Детали мечты
POST /api/v1/donations/pay       # Пожертвовать

GET  /api/v1/news                # Новости
GET  /api/v1/reports             # Отчеты
```

## Технологии

- Go 1.24 + Gin
- PostgreSQL + GORM
- JWT аутентификация
- Docker Compose (локально)

## Файлы для работы

| Файл | Назначение |
|------|------------|
| `openapi.yaml` | Спецификация API - импортируй в Postman или Swagger Editor |
| `FRONTEND.md` | Быстрый старт для фронтендера |
| `API_DOCUMENTATION.md` | Полная документация с примерами |
| `cmd/seed/seed.go` | Мок-данные для тестирования |

---

**Для диплома**: Код прокомментирован на русском, все обработчики реализованы согласно `openapi.yaml`. Фронтенд может использовать API локально.

Репозиторий: https://github.com/stalxr/consellation-react-backend-gogingorn.git
