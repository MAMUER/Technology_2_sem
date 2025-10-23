# Практическая работа №8
# Николаенко Михаил ЭФМО-02-21

## Описание проекта и требования

MongoDB как документная NoSQL-база предоставляет гибкость хранения JSON-подобных документов без жесткой схемы. Официальный драйвер для Go позволяет выполнять CRUD-операции, текстовый поиск и агрегации. Идеально для проектов с изменяющейся структурой данных и быстрым прототипированием.

Для работы с командой make в PowerShell необходимо установить менеджер пакетов Chocolatey и установить команду make

Проект на языке Go (необходима версия 1.21 и выше) с REST-API:

Основные эндпоинты:

- `GET /health` – базовая проверка здоровья сервиса.
- `POST /api/v1/notes` – создание новой заметки.
- `GET /api/v1/notes` – получение списка заметок с пагинацией.
- `GET /api/v1/notes/{id}` – получение заметки по ID.
- `PATCH /api/v1/notes/{id}` – частичное обновление заметки.
- `DELETE /api/v1/notes/{id}` – удаление заметки.
- `GET /api/v1/notes?q={query}` – поиск заметок по заголовку (regex).
- `GET /api/v1/notes/search/text?q={query}` – полнотекстовый поиск по содержимому.
- `GET /api/v1/notes/stats` – общая статистика по заметкам.
- `GET /api/v1/notes/stats/daily?days={n}` – статистика по дням.

## Необходимые пароли

Пользователь сервера
- логин: teacher
- пароль: 1
- IP: 193.233.175.221

Пользователь MongoDB
- логин: teacher
- пароль: teacher_password_123
- порт: 27017
- база данных: pz8

## Команды запуска/сборки

### Запуск MongoDB на сервере:
ssh teacher@193.233.175.221

cd ~/pz8-mongo

docker-compose up -d

### Проверка подключения к MongoDB:
docker-compose exec mongo mongosh -u root -p secret --authenticationDatabase admin --eval "db.runCommand({ping: 1})"

### Запуск тоннеля подключения к серверу (в отдельной консоли):

ssh -L 27017:localhost:27017 teacher@193.233.175.221 -N -o ServerAliveInterval=30

### Сборка приложения:

make build

### Запуск приложения:

make run

### Остановка тоннеля подключения:

make tunnel-stop

### Иснтрукция подключения:

make setup-teacher

### Показать текущие туннели:

make tunnel-status

## Команды:

### Базовая проверка здоровья
curl http://localhost:8080/health

Ответ:

{"status":"ok"}

### Создание заметки:
curl -X POST http://localhost:8080/api/v1/notes ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Первая заметка\",\"content\":\"Текст заметки...\"}"

Ответ:

{"id":"68f62982646a61dc68d7c292","title":"Первая заметка","content":"Текст заметки...","createdAt":"2025-10-20T15:22:26.7451612+03:00","updatedAt":"2025-10-20T15:22:26.7451612+03:00"}

### Получение списка заметок:
#### простой список
curl "http://localhost:8080/api/v1/notes?limit=5&skip=0"

Ответ:

{"notes":[{"id":"68f62982646a61dc68d7c292","title":"Первая заметка","content":"Текст заметки...","createdAt":"2025-10-20T12:22:26.745Z","updatedAt":"2025-10-20T12:22:26.745Z","score":1},{"id":"68f4fe160b08ca63d9be709e","title":"Remote note","content":"Working with remote MongoDB!","createdAt":"2025-10-19T15:04:54.171Z","updatedAt":"2025-10-19T15:04:54.171Z","score":1}],"query":"","searchType":""}

#### поиск по заголовку (regex)
curl "http://localhost:8080/api/v1/notes?q=заметка&limit=5"

Ответ:

{"notes":[{"id":"68f62982646a61dc68d7c292","title":"Первая заметка","content":"Текст заметки...","createdAt":"2025-10-20T12:22:26.745Z","updatedAt":"2025-10-20T12:22:26.745Z","score":1}],"query":"заметка","searchType":""}

#### полнотекстовый поиск
curl "http://localhost:8080/api/v1/notes/search/text?q=программирование&limit=10"

Ответ:

{"notes":null,"query":"программирование","total":0}

### Получение заметки по ID:
curl http://localhost:8080/api/v1/notes/<object_id_here>

Ответ:

{"id":"68f62982646a61dc68d7c292","title":"Первая заметка","content":"Текст заметки...","createdAt":"2025-10-20T12:22:26.745Z","updatedAt":"2025-10-20T12:22:26.745Z"}

### Частичное обновление заметки:
curl -X PATCH http://localhost:8080/api/v1/notes/68f62982646a61dc68d7c292 ^
  -H "Content-Type: application/json" ^
  -d "{\"content\":\"Обновленный текст\"}"

Ответ:

{"id":"68f62982646a61dc68d7c292","title":"Первая заметка","content":"Обновленный текст","createdAt":"2025-10-20T12:22:26.745Z","updatedAt":"2025-10-20T12:25:51.365Z"}

### Удаление заметки:
curl -X DELETE http://localhost:8080/api/v1/notes/<object_id_here>

### Общая статистика:
curl http://localhost:8080/api/v1/notes/stats

Ответ:

{"totalNotes":1,"avgContentLength":28,"maxContentLength":28,"minContentLength":28}

### Статистика по дням:
curl "http://localhost:8080/api/v1/notes/stats/daily?days=7"

Ответ:

{"period":7,"stats":[{"date":"2025-10-19","count":1}]}

## Структура проекта
```
C:.
└───pz8-mongo
    ├───.env.example
    ├───go.mod
    ├───go.sum
    ├───Makefile
    ├───README.md
    │
    ├───bin
    │   └───server.exe
    │
    ├───cmd
    │   └───api
    │       └───main.go
    │
    ├───internal
    │   ├───db
    │   │   └───mongo.go
    │   │
    │   └───notes
    │       ├───handler.go
    │       ├───model.go
    │       ├───repo.go
    │       └───repo_test.go
    │
    └───PR8
```
## Примечания по конфигурации

Подключение к MongoDB происходит через файл .env.example

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR8/Screenshot_1.png)

Создание конфигурационного файла MongoDB на сервере

![фото18](./PR8/Screenshot_18.png)

Установка Docker на сервере

![фото5](./PR8/Screenshot_5.png)

![фото3](./PR8/Screenshot_3.png)

![фото4](./PR8/Screenshot_4.png)

Запуск MongoDB на сервере

![фото6](./PR8/Screenshot_6.png)

![фото14](./PR8/Screenshot_14.png)

Проверка работы MongoDB

![фото7](./PR8/Screenshot_7.png)

Создание пользователя для MongoDB

![фото15](./PR8/Screenshot_15.png)

Подключение к серверу по SSH тоннелю

![фото8](./PR8/Screenshot_8.png)

Проверка и запуск локального приложения

![фото9](./PR8/Screenshot_9.png)

Базовая проверка здоровья
curl http://localhost:8080/health

Создание заметки, простой список заметок, поиск заметок по заголовку (regex), полнотекстовый поиск заметок, получение заметки по ID, частичное обновление заметки, удаление заметки, общая статистика, статистика по дням:

![фото10](./PR8/Screenshot_10.png)

![фото11](./PR8/Screenshot_11.png)

![фото12](./PR8/Screenshot_12.png)

![фото16](./PR8/Screenshot_16.png)

![фото17](./PR8/Screenshot_17.png)

Тест

![фото13](./PR8/Screenshot_13.png)

Структура проекта

![фото19](./PR8/Screenshot_19.png)