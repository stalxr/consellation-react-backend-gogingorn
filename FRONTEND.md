# Гайд для фронтендера

## Что такое seed.go?

`seed.go` — это **Go-программа** (скрипт), которая при запуске:
1. Подключается к PostgreSQL базе данных (той же, что использует API)
2. Создаёт тестовые данные: пользователей, детей с онкологией, пожертвования, новости

Это **не отдельная база данных**, это просто автоматическое заполнение существующей PostgreSQL.

## Быстрый старт

### 1. Запуск бэкенда локально

```bash
# 1. Клонируй репо
git clone https://github.com/stalxr/consellation-react-backend-gogingorn.git
cd consellation-react-backend-gogingorn

# 2. Запусти PostgreSQL
docker-compose up -d

# 3. Запусти API (в одном терминале)
go run cmd/api/main.go

# 4. Заполни мок-данными (в другом терминале)
go run cmd/seed/seed.go
```

API будет на `http://localhost:8080/api/v1`

### 2. Что создаёт seed.go

**Тестовые аккаунты:**
- Админ: `admin@nastenka.ru` / `admin123`
- Юзер: `user@example.com` / `user123`

**Дети с онкологией (10 историй):**
- Маша, 7 лет — лейкемия (нужна химиотерапия)
- Ваня, 5 лет — опухоль мозга (операция в Израиле)
- Катя, 12 лет — рак глаза (протонная терапия в Германии) ✅ собрано
- Дима, 9 лет — саркома (таргетная терапия)
- Настя, 6 лет — реабилитация после нейробластомы
- Лёша, 14 лет — глиома (гамма-нож)
- Полина, 11 лет — CAR-T терапия (рецидив лейкемии)
- Саша, 3 года — трансплантация печени ✅ собрано
- Миша, 8 лет — диагностика лимфомы
- Варя, 16 лет — паллиативная помощь

**Новости и отчёты** — для тестирования разделов

## Варианты работы с мок-данными

### Вариант А: API + База (рекомендуется)
Фронтендер запускает у себя:
```bash
# Терминал 1: API
go run cmd/api/main.go

# Терминал 2: Seed (один раз)
go run cmd/seed/seed.go
```
Потом делает запросы к `http://localhost:8080/api/v1`

### Вариант Б: JSON Mock Server
Если фронтендер использует JSON Server или похожий инструмент:
- Данные можно взять из `openapi.yaml` (там есть примеры ответов)
- Или создать `mock.json` — скажи если нужно

### Вариант В: SQLite / JSON файл
Если нужен просто файл с данными без PostgreSQL:
- Могу создать `mock_data.json` со всеми данными
- Фронтендер читает напрямую или импортирует в свой mock

**Какой вариант удобнее для твоего фронтендера?** Скажи — подготовлю.

## Тестовые запросы

```bash
# Получить все мечты (дети с онкологией)
curl http://localhost:8080/api/v1/dreams

# Войти как админ
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@nastenka.ru","password":"admin123"}'

# Создать пожертвование (нужен токен)
curl -X POST http://localhost:8080/api/v1/donations/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"amount":5000,"dream_id":"11111111-1111-1111-1111-111111111111"}'
```

## Swagger / OpenAPI

Файл `openapi.yaml` содержит полную спецификацию:
- Импортируй в Postman (File → Import)
- Или открой в Swagger Editor: https://editor.swagger.io/

## Документация

- `openapi.yaml` — API спецификация
- `API_DOCUMENTATION.md` — подробное описание эндпоинтов
- `Dream_API_Postman_Collection.json` — готовая коллекция Postman

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
