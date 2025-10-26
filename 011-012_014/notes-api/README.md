# Практическая работа №11-12, 14
# Николаенко Михаил ЭФМО-02-21

## Описание проекта и требования

Освоение принципов проектирования REST API и слоистой архитектуры через реализацию CRUD-интерфейса для системы заметок с использованием code-first подхода, при котором кодовая база предшествует формальному описанию API.

Проект на языке Go (необходима версия 1.21 и выше) с REST-API

- `POST /api/v1/notes` - Создание заметки
- `GET /api/v1/notes` - Получение списка заметок с пагинацией и поиском
- `GET /api/v1/notes/{id}` - Получение конкретной заметки по ID
- `PATCH /api/v1/notes/{id}` - Частичное обновление заметки
- `DELETE /api/v1/notes/{id}` - Удаление заметки

## Необходимые пароли

Подключение к серверу
- логин: teacher
- пароль: 1
- IP: 193.233.175.221

Пользователь PostgreSQL
- логин: teacher_app 
- пароль: secure_password_123
- порт: 5432
- БД: notes

## Команды запуска/сборки

### Подключение к серверу 
ssh teacher@193.233.175.221

### Запуск тоннеля подключения к серверу (в отдельной консоли):

ssh -L 5433:localhost:5432 teacher@193.233.175.221 -N -o ServerAliveInterval=30

### Сборка приложения:

make build

### Запуск приложения:

make run

### Запуск swagger:

make swagger

## Команды:

### Создание заметки
curl -X POST http://localhost:8080/api/v1/notes ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer test-token-12345" ^
-d "{\"title\":\"Первая заметка\", \"content\":\"Это тест\"}"

Ответ:

#### In-memory
{"ID":1,"Title":"Первая заметка","Content":"Это тест","CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":null}

#### PostgreSQL
{"id":301}

### Получение списка заметок
curl -X GET "http://localhost:8080/api/v1/notes?page=1&limit=10&q=заметка" ^
-H "Authorization: Bearer test-token-12345" ^
-H "Content-Type: application/json"

Ответ:

#### In-memory
[{"ID":1,"Title":"Первая заметка","Content":"Это тест","CreatedAt":"2024-01-15T10:30:00Z","UpdatedAt":null}]

#### PostgreSQL
[{"id":301,"title":"Первая заметка","content":"Это тест","created_at":"2025-10-26T19:01:08.493454+03:00"},{"id":300,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":299,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":297,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":296,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":295,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":294,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":293,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":292,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.826133+03:00"},{"id":298,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:52.825408+03:00"}]

### Получение конкретной заметки
curl -X GET http://localhost:8080/api/v1/notes/1 ^
-H "Authorization: Bearer test-token-12345" ^
-H "Content-Type: application/json"

Ответ:

#### In-memory
{"ID":1,"Title":"Первая заметка","Content":"Это тест","CreatedAt":"2024-01-15T10:30:00Z","UpdatedAt":null}

#### PostgreSQL
{"id":1,"title":"Test","content":"Content","created_at":"2025-10-26T18:50:49.029426+03:00"}

### Частичное обновление заметки
curl -X PATCH http://localhost:8080/api/v1/notes/1 ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer test-token-12345" ^
-d "{\"title\":\"Обновленная заметка\"}"

Ответ:

#### In-memory
{"ID":1,"Title":"Обновленная заметка","Content":"Это тест","CreatedAt":"2024-01-15T10:30:00Z","UpdatedAt":"2024-01-15T11:00:00Z"}

#### PostgreSQL
{"id":1,"title":"Обновленная заметка","content":"Content","created_at":"2025-10-26T18:50:49.029426+03:00"}

### Удаление заметки
curl -X DELETE http://localhost:8080/api/v1/notes/1 ^
-H "Authorization: Bearer test-token-12345" ^
-H "Content-Type: application/json"

## Структура проекта
```
C:.
│   .env
│   go.mod
│   go.sum
│   Makefile
│   README.md
│
├───api
│       openapi.yaml
│
├───bin
│       server.exe
│
├───cmd
│   └───api
│           main.go
│
├───docs
│       docs.go
│       swagger.json
│       swagger.yaml
│
├───internal
│   ├───config
│   │       database.go
│   │
│   ├───core
│   │   │   note.go
│   │   │
│   │   └───service
│   │           note_service.go
│   │           transaction_service.go
│   │
│   ├───http
│   │   │   router.go
│   │   │
│   │   └───handlers
│   │           notes.go
│   │
│   └───repo
│           note_mem.go
│           note_postgres.go
│
└───PR11-12_14
```

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR11-12_14/Screenshot_1.png)

![фото5](./PR11-12_14/Screenshot_5.png)

![фото7](./PR11-12_14/Screenshot_7.png)

![фото11](./PR11-12_14/Screenshot_11.png)

Создание БД

![фото10](./PR11-12_14/Screenshot_10.png)

![фото13](./PR11-12_14/Screenshot_13.png)

Добавление прав доступа

![фото14](./PR11-12_14/Screenshot_14.png)

![фото18](./PR11-12_14/Screenshot_18.png)

Добавление shared_preload_libraries

![фото17](./PR11-12_14/Screenshot_17.png)

![фото16](./PR11-12_14/Screenshot_16.png)

Создание файлов для подключения к БД

![фото12](./PR11-12_14/Screenshot_12.png)

Запуск проекта

![фото8](./PR11-12_14/Screenshot_8.png)

Генерация swagger-доков

![фото6](./PR11-12_14/Screenshot_6.png)

Проверка и запуск приложения

![фото2](./PR11-12_14/Screenshot_2.png)

Создание заметки

![фото3](./PR11-12_14/Screenshot_3.png)

Проверка обновленных команд

![фото9](./PR11-12_14/Screenshot_9.png)

EXPLAIN/ANALYZE проблемных запросов

![фото15](./PR11-12_14/Screenshot_15.png)

Статистика запросов из БД

![фото19](./PR11-12_14/Screenshot_19.png)

Нагрузочное тестирование

![фото20](./PR11-12_14/Screenshot_20.png)

![фото21](./PR11-12_14/Screenshot_21.png)

![фото22](./PR11-12_14/Screenshot_22.png)

![фото23](./PR11-12_14/Screenshot_23.png)

Тестирование разных размеров пула

![фото24](./PR11-12_14/Screenshot_24.png)

![фото25](./PR11-12_14/Screenshot_25.png)

![фото26](./PR11-12_14/Screenshot_26.png)

Мониторинг БД в реальном времени

![фото27](./PR11-12_14/Screenshot_27.png)

Проверка обновленных команд

![фото28](./PR11-12_14/Screenshot_28.png)

![фото29](./PR11-12_14/Screenshot_29.png)

Структура проекта

![фото4](./PR11-12_14/Screenshot_4.png)