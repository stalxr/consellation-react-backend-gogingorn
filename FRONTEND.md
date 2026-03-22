# Гайд для фронтендера

## Быстрый старт

API запускается локально на `http://localhost:8080/api/v1`

### 1. Запуск бэкенда локально

```bash
# 1. Клонируй репо
git clone https://github.com/stalxr/consellation-react-backend-gogingorn.git
cd consellation-react-backend-gogingorn

# 2. Запусти PostgreSQL
docker-compose up -d

# 3. Запусти API
go run cmd/api/main.go
```

API будет на `http://localhost:8080/api/v1`

### 2. Тестовые данные

```bash
# Заполни базу мок-данными
go run cmd/seed/seed.go
```

Создаст:
- Админ: `admin@nastenka.ru` / `admin123`
- Юзер: `user@example.com` / `user123`
- 6 мечт, 3 новости, донаты, отчёты

### 3. Swagger / OpenAPI

Файл спецификации: `openapi.yaml`

Можно открыть в [Swagger Editor](https://editor.swagger.io/) или Postman (Import → File → openapi.yaml)

## Документация

| Файл | Описание |
|------|----------|
| `openapi.yaml` | Полная спецификация API (OpenAPI 3.0) |
| `API_DOCUMENTATION.md` | Подробное описание всех эндпоинтов с примерами |
| `Dream_API_Postman_Collection.json` | Коллекция Postman для тестирования |

## Ключевые эндпоинты

```
POST /api/v1/auth/login          # Вход
POST /api/v1/auth/register       # Регистрация
GET  /api/v1/auth/me             # Профиль (нужен токен)

GET  /api/v1/dreams              # Список мечт
GET  /api/v1/dreams/{id}         # Детали мечты
POST /api/v1/donations/pay       # Создать пожертвование

GET  /api/v1/news                # Новости
GET  /api/v1/reports             # Отчёты
```

## CORS

Бэкенд разрешён для всех origin:
```javascript
fetch('http://localhost:8080/api/v1/dreams')
// или с proxy в React
```

## Аутентификация

```javascript
// Сохраняй токен после логина
localStorage.setItem('token', data.access_token)

// Используй в запросах
fetch('/api/v1/auth/me', {
  headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
})
```

## Структура проекта (для понимания)

```
cmd/api/main.go           # Точка входа
cmd/seed/seed.go          # Мок-данные

internal/
├── handlers/             # HTTP обработчики
├── services/             # Бизнес-логика
├── repository/           # База данных
├── middleware/           # JWT, CORS
└── models/               # Структуры данных

openapi.yaml              # Спецификация API
```

## Если что-то не работает

1. Проверь что PostgreSQL запущен: `docker ps`
2. Проверь логи API: `go run cmd/api/main.go` смотри вывод
3. Проверь `.env` файл (JWT_SECRET должен быть)

---

Всё что нужно для интеграции — в `openapi.yaml`. Swagger Editor покажет все схемы запросов/ответов.
