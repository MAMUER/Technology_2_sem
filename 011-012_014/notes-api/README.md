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

## Команды запуска/сборки

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

{"ID":1,"Title":"Первая заметка","Content":"Это тест","CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":null}

### Получение списка заметок
curl -X GET "http://localhost:8080/api/v1/notes?page=1&limit=10&q=заметка" ^
-H "Authorization: Bearer test-token-12345" ^
-H "Content-Type: application/json"

Ответ:

[{"ID":1,"Title":"Первая заметка","Content":"Это тест","CreatedAt":"2024-01-15T10:30:00Z","UpdatedAt":null}]

### Получение конкретной заметки
curl -X GET http://localhost:8080/api/v1/notes/1 ^
-H "Authorization: Bearer test-token-12345" ^
-H "Content-Type: application/json"

Ответ:

{"ID":1,"Title":"Первая заметка","Content":"Это тест","CreatedAt":"2024-01-15T10:30:00Z","UpdatedAt":null}

### Частичное обновление заметки
curl -X PATCH http://localhost:8080/api/v1/notes/1 ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer test-token-12345" ^
-d "{\"title\":\"Обновленная заметка\"}"

Ответ:

{"ID":1,"Title":"Обновленная заметка","Content":"Это тест","CreatedAt":"2024-01-15T10:30:00Z","UpdatedAt":"2024-01-15T11:00:00Z"}

### Удаление заметки
curl -X DELETE http://localhost:8080/api/v1/notes/1 ^
-H "Authorization: Bearer test-token-12345" ^
-H "Content-Type: application/json"

## Структура проекта
```
C:.
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
│   ├───core
│   │   │   note.go
│   │   │
│   │   └───service
│   │           note_service.go
│   │
│   ├───http
│   │   │   router.go
│   │   │
│   │   └───handlers
│   │           notes.go
│   │
│   └───repo
│           note_mem.go
│
└───PR11-12
```

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR11-12/Screenshot_1.png)

![фото5](./PR11-12/Screenshot_5.png)

![фото7](./PR11-12/Screenshot_7.png)

Запуск проекта

![фото8](./PR11-12/Screenshot_8.png)

Генерация swagger-доков

![фото6](./PR11-12/Screenshot_6.png)

Проверка и запуск приложения

![фото2](./PR11-12/Screenshot_2.png)

Создание заметки

![фото3](./PR11-12/Screenshot_3.png)

Проверка обновленных команд

![фото9](./PR11-12/Screenshot_9.png)

Структура проекта

![фото4](./PR11-12/Screenshot_4.png)