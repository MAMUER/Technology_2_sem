# Практическая работа №4
# Николаенко Михаил ЭФМО-02-21

## Описание проекта и требования

### Требования

Для работы с командой make в PowerShell необходимо установить менеджер пакетов Chocolatey и установить команду make.

Проект представляет собой простой HTTP-сервер на языке Go (требуется версия 1.21 и выше) с REST-API.:

Основные эндпоинты:

- `GET /health` – проверка состояния сервера.
- `GET /api/v1/tasks` – получение списка всех задач.
- `POST /api/v1/tasks` – создание новой задачи.
- `GET /api/v1/tasks/{id}` – получение конкретной задачи по ID.
- `PATCH /api/v1/tasks/{id}` – отметить задачу выполненной.
- `GET /api/v1/tasks?q=TEXT` – поиск задач с фильтром.
- `DELETE /api/v1/tasks/{id}` – удалить задачу.

## Команды запуска и сборки

### Сборка приложения:

make build

### Запуск приложения:

make run

### Примеры запросов:

#### Проверка состояния сервера:

curl http://localhost:8080/health

Ответ:

{
  "status": "ok"
}

#### Получение списка задач с пагинацией и фильтром:

curl "http://localhost:8080/api/v1/tasks?page=1&limit=10&done=false"

#### Получение списка задач:

curl http://localhost:8080/api/v1/tasks

Ответ:

[{
  "id":1,"title":"TEXT","done":false},
  {"id":2,"title":"TEXT","done":false},
...}]

#### Создание новой задачи:

curl -Method POST http://localhost:8080/api/v1/tasks `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"title":"TEXT"}'

Ответ:

{
  "id":1,
  "title":"TEXT",
  "done":false
}

#### Обновление задачи:

curl -Method PUT http://localhost:8080/api/v1/tasks/1 `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"title":"NEWTEXT","done":true}'


#### Получение задачи по ID:

curl http://localhost:8080/api/v1/tasks/1

Ответ:

{
  "id":1,
  "title":"TEXT",
  "done":false
}

#### Отметить задачу выполненной:

curl http://localhost:8080/api/v1/tasks/1 -Method PATCH

Ответ:

{
  "id":1,
  "title":"TEXT",
  "done":true
}

#### Поиск задач с фильтром:

curl http://localhost:8080/api/v1/tasks?q=TEXT

Ответ:

[
  {
    "id":1,
    "title":"TEXT",
    "done":false
  }
]

#### Удалить задачу:

curl http://localhost:8080/api/v1/tasks/1 -Method DELETE

## Структура проекта
```
C:.
└───pz4-todo
    ├───go.mod
    ├───go.sum
    ├───main.go
    ├───Makefile
    ├───README.md
    ├───tasks.json
    │
    ├───bin
    │   └───server.exe
    │
    ├───internal
    │   └───task
    │       ├───handler.go
    │       ├───model.go
    │       └───repo.go
    │
    ├───pkg
    │   └───middleware
    │       ├───cors.go
    │       └───logger.go
    │
    └───PR4
```

## Примечания по конфигурации

- Сервер использует хранение данных в памяти с сохранением в файл (in-memory storage + persistence).
- Логи всех входящих запросов видны в консоли.
- По умолчанию сервер слушает порт 8080.
- Можно задать переменную окружения PORT для изменения порта запуска.
- В проекте используется middleware для логирования запросов и CORS.

## Скриншоты работы проекта

Инициализация модуля и установка зависимостей

![фото1](./PR4/Screenshot_1.png)

Запуск и логи проекта

![фото2](./PR4/Screenshot_8.png)

Проверка в браузере (/health)

![фото3](./PR4/Screenshot_2.png)

Создание задачи через curl (/tasks -POST)

![фото4](./PR4/Screenshot_3.png)

Проверка через curl (/tasks)

![фото5](./PR4/Screenshot_4.png)

Проверка через curl (/tasks/{id})

![фото6](./PR4/Screenshot_5.png)

Обновление задачи через curl (/tasks/{id} -PUT)

![фото7](./PR4/Screenshot_6.png)

Удаление задачи через curl (/tasks/{id} -DELETE)

![фото8](./PR4/Screenshot_7.png)

Проверка через curl (/tasks?page=1&limit=10)

![фото10](./PR4/Screenshot_13.png)

Проверка через curl (/tasks?done=true)

![фото9](./PR4/Screenshot_14.png)

Создание task (новый api/v1)

![фото11](./PR4/Screenshot_11.png)

Вывод всех tasks (новый api/v1)

![фото12](./PR4/Screenshot_12.png)

Проверки форматирования кода и базовая проверка

![фото14](./PR3/Screenshot_9.png)

Структура проекта

![фото13](./PR4/Screenshot_15.png)