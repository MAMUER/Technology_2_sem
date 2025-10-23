# Практическая работа №8
# Николаенко Михаил ЭФМО-02-21

## Описание проекта и требования

REST API с CRUD-функциональностью для управления заметками.

Проект на языке Go (необходима версия 1.21 и выше) с REST-API

- `POST /api/v1/notes` - Создание заметки

## Необходимые пароли

## Команды запуска/сборки

### Сборка приложения:

make build

### Запуск приложения:

make run

### Инструкция:

make help

## Команды:

### Создание заметки
curl -X POST http://localhost:8080/api/v1/notes ^
-H "Content-Type: application/json" ^
-d "{\"title\":\"Первая заметка\", \"content\":\"Это тест\"}"

Ответ:

{"ID":1,"Title":"Первая заметка","Content":"Это тест","CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":null}

## Структура проекта
```
C:.
├───go.mod
├───go.sum
├───Makefile
├───README.md
│
├───api
│   └───openapi.yaml
│
├───bin
│   └───server.exe
│
├───cmd
│   └───api
│       └───main.go
│
├───internal
│   ├───core
│   │   ├───note.go
│   │   │
│   │   └───service
│   │       └───note_service.go
│   │
│   ├───http
│   │   ├───router.go
│   │   │
│   │   └───handlers
│   │       └───notes.go
│   │
│   └───repo
│       └───note_mem.go
│
└───PR11
```

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR11/Screenshot_1.png)

Проверка и запуск приложения

![фото2](./PR11/Screenshot_2.png)

Создание заметки

![фото3](./PR11/Screenshot_3.png)

Структура проекта

![фото4](./PR11/Screenshot_4.png)