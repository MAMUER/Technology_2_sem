# Практическая работа №3
# Николаенко Михаил ЭФМО-02-21

## Описание проекта и требования

### Требования

Проект представляет собой простой HTTP-сервер на языке Go (необходима версия 1.21 и выше) с REST-API:

Основные эндпоинты:

- `GET /health` – проверка состояния сервера.
- `GET /tasks` – получение списка всех задач.
- `POST /tasks` – создание новой задачи.
- `GET /tasks/{id}` – получение конкретной задачи по ID.

## Команды запуска и сборки

### Сборка приложения:

go build -o bin\server.exe ./cmd/server

### Запуск приложения:

.\bin\server.exe

### Примеры запросов:

#### Проверка состояния сервера:

curl http://localhost:8080/health

Ответ:

{
  "status": "ok"
}

#### Получение списка задач:

curl http://localhost:8080/tasks

Ответ:

[{
  "id":1,"title":"TEXT","done":false},
  {"id":2,"title":"TEXT","done":false},
...}]

#### Создание новой задачи:

curl -Method POST http://localhost:8080/tasks `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"title":"TEXT"}'

Ответ:

{
  "id":1,"title":"TEXT","done":false
}

#### Получение задачи по ID:

curl http://localhost:8080/tasks/1

Ответ:

{
  "id":1,"title":"TEXT","done":false
}

## Структура проекта
```
C:.
└───pz3-http
    ├───go.mod
    ├───README.md
    │
    ├───bin
    │   ├───http.exe
    │   └───server.exe
    │
    ├───cmd
    │   └───server
    │       └───main.go
    │
    ├───internal
    │   ├───api
    │   │   ├───handlers.go
    │   │   ├───middleware.go
    │   │   └───responses.go
    │   │
    │   └───storage
    │       └───memory.go
    │
    └───PR3
```

## Примечания по конфигурации

- Сервер использует память для хранения данных (in-memory storage) и логирует все входящие запросы.

- По умолчанию сервер слушает порт 8080.

- Порт можно изменить в параметре http.ListenAndServe(":8080", handler) в main.go.

- Используется middleware для логирования запросов


## Скриншоты работы проекта

Проверка наличия ПО

![фото1](./PR3/Screenshot_4.png)

Инициализация проекта

![фото2](./PR3/Screenshot_4.png)

Запуск сервера и логи во время работы

![фото3](./PR3/Screenshot_13.png)

Проверка работы ручек в браузере (вывод текста, /ping, /fail)

![фото4](./PR3/Screenshot_1.png)
![фото5](./PR3/Screenshot_2.png)
![фото5](./PR3/Screenshot_14.png)

Проверка через curl (/)

![фото6](./PR3/Screenshot_8.png)

Проверка через curl (/ping)

![фото7](./PR3/Screenshot_7.png)

Проверка через curl (/fail)

![фото8](./PR3/Screenshot_12.png)

Сборка бинарника и его запуск

![фото9](./PR3/Screenshot_6.png)

Проверка через curl (Запросы через GitBush)

![фото10](./PR3/Screenshot_9.png)
![фото11](./PR3/Screenshot_10.png)

Структура проекта

![фото12](./PR3/Screenshot_11.png)