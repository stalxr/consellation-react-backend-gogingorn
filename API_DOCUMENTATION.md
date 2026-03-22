# API Documentation for Frontend Developer

## Базовый URL
```
http://localhost:8080/api/v1
```

## Swagger / OpenAPI
```
http://localhost:8080/swagger/index.html
```

---

## 🔐 Аутентификация

### Регистрация (донор)
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "Иван Иванов"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Вход
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Получить профиль (требуется JWT)
```http
GET /auth/me
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "full_name": "Иван Иванов",
  "avatar_url": "https://...",
  "role": "user",
  "total_donated": 1500
}
```

---

## 💫 Мечты (Dreams)

### Список всех мечт (публичный)
```http
GET /dreams?status=active&page=1&limit=9
```

**Query параметры:**
- `status` (optional): `active` | `completed`
- `page` (optional): номер страницы, default: 1
- `limit` (optional): количество на странице, default: 9

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Мечта Вани: увидеть море",
      "slug": "vanya-more",
      "short_description": "...",
      "full_description": "<p>...</p>",
      "target_amount": 50000,
      "collected_amount": 12500,
      "status": "active",
      "cover_image": "/uploads/images/...",
      "gallery_images": [],
      "created_at": "2024-01-01T00:00:00Z",
      "closed_at": null
    }
  ],
  "total": 25
}
```

### Детальная информация о мечте
```http
GET /dreams/{id}
```

### Создать мечту (только админ)
```http
POST /dreams
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "title": "Новая мечта",
  "short_description": "Описание",
  "full_description": "<p>Полное описание</p>",
  "target_amount": 100000,
  "cover_image": "/uploads/images/photo.jpg"
}
```

### Редактировать мечту (только админ)
```http
PUT /dreams/{id}
Authorization: Bearer <admin_token>
Content-Type: application/json
```

### Удалить мечту (только админ)
```http
DELETE /dreams/{id}
Authorization: Bearer <admin_token>
```

### Список доноров мечты
```http
GET /dreams/{id}/donors
```

**Response:**
```json
[
  {
    "amount": 1000,
    "donor_name": "Алексей",
    "created_at": "2024-01-15T10:30:00Z",
    "is_anonymous": false
  }
]
```

---

## 💳 Пожертвования (Donations)

### Создать платеж
```http
POST /donations/pay
Content-Type: application/json

{
  "amount": 1000,
  "dream_id": "uuid",
  "email": "donor@example.com",
  "is_anonymous": false,
  "comment": "Удачи!"
}
```

**Response:**
```json
{
  "payment_url": "https://payment.example.com/pay?donation_id=..."
}
```

### История моих пожертвований (требуется JWT)
```http
GET /donations/my
Authorization: Bearer <access_token>
```

**Response:**
```json
[
  {
    "id": "uuid",
    "amount": 1000,
    "dream_title": "Мечта Вани: увидеть море",
    "date": "2024-01-15T10:30:00Z",
    "status": "completed"
  }
]
```

---

## 📰 Новости (News)

### Список новостей
```http
GET /news
```

### Создать новость (только админ)
```http
POST /news
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "title": "Заголовок новости",
  "preview_text": "Краткое описание...",
  "content": "<p>Полный текст...</p>",
  "image_url": "/uploads/images/news.jpg"
}
```

---

## 📊 Отчеты (Reports)

### Список отчетов
```http
GET /reports
```

### Загрузить отчет PDF (только админ)
```http
POST /reports
Authorization: Bearer <admin_token>
Content-Type: multipart/form-data

file: <PDF файл>
title: "Отчет за январь"
year: 2024
```

---

## 📤 Загрузка файлов (Upload)

### Загрузить изображение (только админ)
```http
POST /upload
Authorization: Bearer <admin_token>
Content-Type: multipart/form-data

file: <изображение jpg/png/gif/webp>
```

**Response:**
```json
{
  "url": "/uploads/images/uuid.jpg"
}
```

---

## 🔔 Webhook для платежей

### Callback от платежной системы
```http
POST /payments/callback
Content-Type: application/json

{
  "donation_id": "uuid",
  "status": "completed",
  "amount": 1000,
  "payment_id": "transaction_123"
}
```

---

## 🔑 Тестовые учетные записи

После запуска `go run ./cmd/seed/seed.go`:

| Роль | Email | Пароль |
|------|-------|--------|
| Админ | `admin@nastenka.ru` | `admin123` |
| Юзер | `user@example.com` | `user123` |

---

## 🚀 Как запустить

```bash
# 1. Запустить PostgreSQL и API
docker-compose up -d

# 2. Заполнить тестовыми данными
go run ./cmd/seed/seed.go

# 3. Или запустить API локально
cd cmd/api
go run main.go
```

API будет доступно по адресу: **http://localhost:8080**

Swagger UI: **http://localhost:8080/swagger/index.html**

---

## 📄 OpenAPI Specification

Файл спецификации: `openapi.yaml` в корне проекта.

Можно импортировать в:
- Swagger Editor (https://editor.swagger.io/)
- Postman
- SwaggerHub
- Insomnia
