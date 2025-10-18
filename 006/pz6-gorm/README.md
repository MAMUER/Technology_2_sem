# Практическая работа №6
# Николаенко Михаил ЭФМО-02-25

## Описание проекта и требования

GORM как ORM устраняет рутину работы с SQL, позволяя описывать модели Go-структурами с автоматической генерацией таблиц. Миграции и связи из коробки экономят время, а безопасность запросов защищает от инъекций. Идеально для быстрого старта и учебных проектов.

## Необходимые пароли

Пользователь сервера
логин: teacher
пароль: 1

Пользователь БД
логин: teacher_app 
пароль: secure_password_123

## Команды запуска/сборки

### Сборка приложения:

make build

### Запуск приложения:

make run

### Запуск тоннеля подключения к серверу (в отдельной консоли):

ssh -L 5433:localhost:5432 teacher@193.233.175.221 -N -o ServerAliveInterval=30

### Остановка тоннеля подключения:

make tunnel-stop

### Проверка подключения:

make check-db

### Иснтрукция подключения:

make setup-teacher

### Показать текущие туннели:

make tunnel-status

## Команды:

# здоровье
curl http://localhost:8080/health

# создаём пользователя
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d "{\"name\":\"Alice\",\"email\":\"alice@example.com\"}"

# создаём заметку с тегами
curl -X POST http://localhost:8080/notes -H "Content-Type: application/json" -d "{\"title\":\"Первая заметка\",\"content\":\"Текст...\",\"userId\":1,\"tags\":[\"go\",\"gorm\"]}"

# получаем заметку с автором и тегами
curl http://localhost:8080/notes/1

## Структура проекта
```
C:.
└───pz6-gorm
    ├───.env
    ├───go.mod
    ├───go.sum
    ├───Makefile
    ├───README.md
    │
    ├───bin
    │   └───server.exe
    │
    ├───cmd
    │   └───server
    │       └───main.go
    │
    ├───internal
    │   ├───db
    │   │   └───postgres.go
    │   │
    │   ├───httpapi
    │   │   ├───handlers.go
    │   │   └───router.go
    │   │
    │   └───models
    │       └───models.go
    │
    └───PR6
```
## Примечания по конфигурации

- Подключение к PostgreSQL происходит через строку подключения из переменной окружения DB_DSN

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR6/Screenshot_1.png)

Выдача прав пользователю

![фото2](./PR6/Screenshot_2.png)

Запуск проекта

![фото3](./PR6/Screenshot_3.png)

здоровье

![фото4](./PR6/Screenshot_4.png)

создаём пользователя, создаём заметку с тегами, получаем заметку с автором и тегами

![фото4](./PR6/Screenshot_5.png)

Структура проекта

![фото8](./PR6/Screenshot_6.png)