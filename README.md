# Благотворительный платформенный бэкенд - Дипломный проект

## Обзор проекта

Это production-ready бэкенд для благотворительной платформы, разработанный на Go с использованием Clean Architecture. Проект строго следует спецификации OpenAPI 3.0 и готов к интеграции с React фронтендом.

### Основные технологии

- **Go 1.25.1** - основной язык программирования
- **Gin** - веб-фреймворк для HTTP роутинга
- **GORM** - ORM для работы с PostgreSQL
- **PostgreSQL 15** - основная база данных
- **JWT** - аутентификация и авторизация
- **Docker Compose** - контейнеризация базы данных
- **OpenAPI 3.0** - спецификация API

## Архитектура проекта

### Clean Architecture

```
cmd/api/          # Точка входа приложения
internal/
├── handlers/     # HTTP обработчики (controllers)
├── services/     # Бизнес-логика
├── repository/   # Работа с базой данных
├── middleware/   # Middleware (CORS, Auth)
└── models/       # Модели данных
pkg/
├── db/          # Подключение к БД
└── utils/       # Утилиты (JWT, пароли)
```

### Слои архитектуры

1. **Handlers** - обрабатывают HTTP запросы, валидация данных
2. **Services** - бизнес-логика, работа с несколькими репозиториями
3. **Repository** - прямой доступ к базе данных через GORM
4. **Models** - структуры данных с GORM тегами

## Установка и запуск

### Требования

- Go 1.25.1 или выше
- Docker и Docker Compose
- Git

### Шаг 1: Клонирование репозитория

```bash
git clone https://github.com/stalxr/consellation-react-backend-gogingorn.git
cd consellation-react-backend-gogingorn
```

### Шаг 2: Настройка переменных окружения

Файл `.env` уже создан с необходимыми настройками:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=charity
DB_PORT=5432
JWT_SECRET=diploma_secret_key_2026
PORT=8080
```

### Шаг 3: Запуск PostgreSQL через Docker

```bash
docker-compose up -d
```

Это запустит PostgreSQL 15 с постоянным хранением данных.

### Шаг 4: Установка зависимостей и запуск

```bash
go mod tidy
go run cmd/api/main.go
```

Сервер запустится на порту 8080:
- API: `http://localhost:8080/api/v1`
- Статические файлы: `http://localhost:8080/uploads`

## API Эндпоинты

### Аутентификация

- `POST /api/v1/auth/register` - Регистрация донора
- `POST /api/v1/auth/login` - Вход (админ или донор)
- `GET /api/v1/auth/me` - Профиль текущего пользователя

### Мечты (Dreams)

- `GET /api/v1/dreams` - Список мечт (с фильтрацией по статусу)
- `GET /api/v1/dreams/{id}` - Детальная информация о мечте
- `GET /api/v1/dreams/{id}/donors` - Список доноров мечты
- `POST /api/v1/dreams` - Создание мечты (только админ)
- `PUT /api/v1/dreams/{id}` - Обновление мечты (только админ)
- `DELETE /api/v1/dreams/{id}` - Удаление мечты (только админ)

### Пожертвования (Donations)

- `POST /api/v1/donations/pay` - Создание платежа (публичный)
- `GET /api/v1/donations/my` - История пожертвований пользователя

### Новости (News)

- `GET /api/v1/news` - Список новостей
- `POST /api/v1/news` - Создание новости (только админ)

### Отчеты (Reports)

- `GET /api/v1/reports` - Финансовые отчеты
- `POST /api/v1/reports` - Загрузка PDF отчета (только админ)

### Загрузка файлов

- `POST /api/v1/upload` - Загрузка изображений (только админ)

## Тестирование API

### 1. Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "full_name": "Иван Иванов"
  }'
```

### 2. Вход в систему

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Сохраните `access_token` из ответа для следующих запросов.

### 3. Получение списка мечт

```bash
curl http://localhost:8080/api/v1/dreams
```

### 4. Создание пожертвования

```bash
curl -X POST http://localhost:8080/api/v1/donations/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "amount": 1000,
    "dream_id": "dream-uuid-here",
    "is_anonymous": false
  }'
```

## Интеграция с React фронтендом

### CORS настройки

Бэкенд настроен для работы с React через `gin-contrib/cors`:

```go
router.Use(cors.New(cors.Config{
    AllowAllOrigins: true,
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
    AllowCredentials: true,
}))
```

### Аутентификация в React

```javascript
// Логин
const login = async (email, password) => {
  const response = await fetch('/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  const data = await response.json();
  localStorage.setItem('token', data.access_token);
};

// Запрос с токеном
const getProfile = async () => {
  const token = localStorage.getItem('token');
  const response = await fetch('/api/v1/auth/me', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return response.json();
};
```

### Загрузка файлов

```javascript
const uploadImage = async (file) => {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('/api/v1/upload', {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` },
    body: formData
  });
  
  return response.json(); // { url: "/uploads/images/filename.jpg" }
};
```

## Структура базы данных

### Таблицы

1. **users** - Пользователи (доноры и администраторы)
2. **dreams** - Мечты (благотворительные сборы)
3. **donations** - Пожертвования
4. **news** - Новости фонда
5. **reports** - Финансовые отчеты

### Связи

- `donations.dream_id` → `dreams.id`
- `donations.user_id` → `users.id` (nullable для анонимных)

## Безопасность

### JWT Аутентификация

- Access токены с коротким сроком действия
- Refresh токены для обновления сессии
- Ролевая модель: "user" и "admin"

### Валидация данных

- Пароли хэшируются через bcrypt
- Валидация email формата
- Проверка обязательных полей
- Ограничение размера загружаемых файлов

## Развертывание (Production)

### Docker

```dockerfile
FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
```

### Environment Variables для Production

```env
GIN_MODE=release
DB_HOST=your-db-host
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=charity_prod
JWT_SECRET=your-super-secret-jwt-key
PORT=8080
```

## Мониторинг и логирование

### Структурированные логи

Приложение использует стандартный logger Go с информативными сообщениями:

```go
log.Println("✅ База данных успешно подключена и миграции выполнены")
log.Printf("🚀 Сервер запущен на порту %s", port)
```

### Health check

```bash
curl http://localhost:8080/api/v1/dreams
```

Если возвращает JSON с мечтами - сервер работает корректно.

## Возможные улучшения

1. **Кэширование** - Redis для частых запросов
2. **Очередь задач** - для обработки платежей
3. **Мониторинг** - Prometheus + Grafana
4. **Тестирование** - unit и integration тесты
5. **CI/CD** - GitHub Actions для автоматического развертывания

## Поддержка и вопросы

Проект полностью документирован русскими комментариями в коде для защиты диплома. Все функции имеют подробные описания на русском языке.

Репозиторий: https://github.com/stalxr/consellation-react-backend-gogingorn.git
