# Практическая работа №3
# Николаенко Михаил ЭФМО-02-25

## Описание проекта и требования

Простой HTTP-сервер для управления задачами (To-Do list) на стандартной библиотеке Go net/http с поддеркой CRUD операций и фильтрацией.

### Требования
- Go версии 1.21 и выше
- Для работы с командой make в PowerShell необходимо установить менеджер пакетов Chocolatey и установить команду make

## Основные эндпоинты
- `GET /health` – проверка состояния сервера
- `GET /tasks` – получение списка всех задач (с поддержкой фильтрации)
- `POST /tasks` – создание новой задачи
- `GET /tasks/{id}` – получение конкретной задачи по ID
- `PATCH /tasks/{id}` – отметить задачу выполненной
- `DELETE /tasks/{id}` – удалить задачу

## Команды запуска и сборки

### Сборка приложения

make build

### Запуск приложения

make run

### Проверка кода и форматирование

make check

### Быстрая сборка и запуск

make fast

### Показать структуру проекта

make tree

### Запуск на определенном порту

make env PORT=####

### Помощь

make help

## Примеры запросов:

### Проверка состояния сервера:

curl http://localhost:8080/health

Ответ:

{
  "status": "ok"
}

### Получение списка задач:

curl http://localhost:8080/tasks

Ответ:

[{
  "id":1,"title":"TEXT","done":false},
  {"id":2,"title":"TEXT","done":false},
...}]

### Создание новой задачи:

curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d "{\"title\":\"TEXT\"}"

Ответ:

{
  "id":1,"title":"TEXT","done":false
}

### Получение задачи по ID:

curl http://localhost:8080/tasks/1

Ответ:

{
  "id":1,"title":"TEXT","done":false
}

### Отметить задачу выполненной:

curl -X PATCH http://localhost:8080/tasks/1

Ответ:

{
  "id":1,"title":"TEXT","done":true
}

### Поиск задач с фильтром:

curl http://localhost:8080/tasks?q=TEXT

Ответ:

{
  "id":1,"title":"TEXT","done":false
}

### Удалить задачу:

curl -X DELETE http://localhost:8080/tasks/1

## Структура проекта
```
C:.
│   .env
│   go.mod
│   go.sum
│   Makefile
│   README.md
│   requests.md
│
├───bin
│       server.exe
│
├───cmd
│   └───server
│           main.go
│
├───internal
│   ├───api
│   │       add.go
│   │       handlers.go
│   │       handlers_test.go
│   │       middleware.go
│   │       responses.go
│   │
│   └───storage
│           memory.go
│
└───PR3
```

## Примечания по конфигурации

- По умолчанию сервер слушает порт 8080.

- Переменная окружения `PORT` задаёт порт для запуска HTTP сервера.

## Скриншоты работы проекта

Проверка наличия ПО (+ установка доп. ПО)

![фото1](./PR3/Screenshot_7.png)

![фото25](./PR3/Screenshot_17.png)

![фото26](./PR3/Screenshot_19.png)

Инициализация проекта (+ сборка и проверка)

![фото2](./PR3/Screenshot_7.png)

![фото27](./PR3/Screenshot_18.png)

Запуск сервера и логи во время работы

![фото3](./PR3/Screenshot_6.png)

Проверка через curl (/health)

![фото4](./PR3/Screenshot_1.png)

Создание задачи через curl (/tasks -POST)

![фото5](./PR3/Screenshot_3.png)

Проверка через curl (/tasks)

![фото6](./PR3/Screenshot_20.png)

Проверка через curl (/tasks/{id})

![фото7](./PR3/Screenshot_8.png)

Проверка через curl (/tasks?q=TEXT)

![фото8](./PR3/Screenshot_21.png)

Проверка через curl (/tasks/{id} -DELETE)

![фото9](./PR3/Screenshot_22.png)

Проверка через curl (/tasks/{id} -PATCH)

![фото10](./PR3/Screenshot_23.png)

Проверка через Invoke-RestMethod

![фото11](./PR3/Screenshot_12.png)

![фото12](./PR3/Screenshot_13.png)

![фото13](./PR3/Screenshot_14.png)

![фото14](./PR3/Screenshot_15.png)

Проверка через curl (Запросы через GitBush)

![фото15](./PR3/Screenshot_9.png)

![фото16](./PR3/Screenshot_5.png)

![фото17](./PR3/Screenshot_10.png)

![фото18](./PR3/Screenshot_11.png)

![фото19](./PR3/Screenshot_16.png)

Проверки форматирования кода и базовая проверка

![фото25](./PR3/Screenshot_25.png)

Структура проекта

![фото20](./PR3/Screenshot_24.png)