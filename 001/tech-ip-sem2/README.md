# Практическое занятие №1
# Разделение монолита на 2 микросервиса. Взаимодействие через HTTP

**Студент:** Николаенко Михаил ЭФМО-02-25

## 1. Описание границ сервисов

### Auth Service
Отвечает за аутентификацию и проверку доступа:
- Выдача токенов по учетным данным (упрощенная модель)
- Валидация токенов для других сервисов
- Хранение информации о пользователях (в памяти)

### Tasks Service
Управляет задачами пользователей:
- CRUD операции с задачами
- Хранение задач в памяти (map[id]Task)
- Проверка доступа через Auth Service перед каждой операцией
- Прокидывание request-id для трассировки запросов

4. Переменные окружения
Auth Service
Переменная	Значение по умолчанию	Описание
AUTH_PORT	8081	Порт для запуска Auth сервиса
Tasks Service
Переменная	Значение по умолчанию	Описание
TASKS_PORT	8082	Порт для запуска Tasks сервиса
AUTH_BASE_URL	http://localhost:8081	Базовый URL Auth сервиса
5. Инструкция по запуску
Предварительные требования
Go версии 1.21 или выше

Make (для удобства)

Git

cd tech-ip-sem2
Запуск Auth Service
bash
# Терминал 1
cd services/auth
export AUTH_PORT=8081
go run ./cmd/auth
Запуск Tasks Service
bash
# Терминал 2
cd services/tasks
export TASKS_PORT=8082
export AUTH_BASE_URL=http://localhost:8081
go run ./cmd/tasks
Запуск с использованием Make (рекомендуется)
bash
# Запуск Auth в отдельном терминале
make run-auth

# Запуск Tasks в отдельном терминале
make run-tasks
Быстрый запуск обоих сервисов
bash
# Терминал 1
make fast-auth

# Терминал 2
make fast-tasks
6. Тестирование через curl
6.1. Получение токена
bash
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: req-001" \
  -d '{"username":"student","password":"student"}'
6.2. Проверка токена напрямую
bash
curl -i http://localhost:8081/v1/auth/verify \
  -H "Authorization: Bearer demo-token-for-student" \
  -H "X-Request-ID: req-002"
6.3. Создание задачи
bash
curl -i -X POST http://localhost:8082/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token-for-student" \
  -H "X-Request-ID: req-003" \
  -d '{"title":"Do PZ17","description":"split services","due_date":"2026-01-10"}'
6.4. Получение списка задач
bash
curl -i http://localhost:8082/v1/tasks \
  -H "Authorization: Bearer demo-token-for-student" \
  -H "X-Request-ID: req-004"
6.5. Попытка без токена (должна вернуть 401)
bash
curl -i http://localhost:8082/v1/tasks \
  -H "X-Request-ID: req-005"
7. Скриншоты с подтверждением request-id
Лог Auth Service с request-id
https://./screenshots/auth_log.png

Лог Tasks Service с request-id
https://./screenshots/tasks_log.png

Пример прокидывания request-id через сервисы
https://./screenshots/request_id_propagation.png

8. Структура проекта
text
tech-ip-sem2/
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── docs/
│   └── pz17_api.md
├── services/
│   ├── auth/
│   │   ├── cmd/
│   │   │   └── auth/
│   │   │       └── main.go
│   │   └── internal/
│   │       ├── http/
│   │       │   └── handlers.go
│   │       └── service/
│   │           └── auth.go
│   └── tasks/
│       ├── cmd/
│       │   └── tasks/
│       │       └── main.go
│       └── internal/
│           ├── client/
│           │   └── authclient/
│           │       └── client.go
│           ├── http/
│           │   └── handlers.go
│           └── service/
│               └── tasks.go
└── shared/
    ├── httpx/
    │   └── client.go
    └── middleware/
        ├── logging.go
        └── requestid.go
9. Makefile команды
Команда	Описание
make run-auth	Запуск Auth сервиса
make run-tasks	Запуск Tasks сервиса
make build	Сборка обоих сервисов
make check	Проверка кода (vet, fmt)
make test	Тестирование API
make curl-examples	Показать примеры curl
make tree	Показать структуру проекта
make help	Показать справку