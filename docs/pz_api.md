# Примечания по конфигурации
## Переменные окружения
| Переменная | Значение по умолчанию | Описание |
|------------|----------------------|----------|
| `GRAPHQL_PORT` | 8090 | Порт сервера |
| `AUTH_PORT` | 8081 | Порт HTTP сервера Auth |
| `AUTH_BASE_URL` | http://193.233.175.221:8081 | Базовый URL Auth сервиса |
| `AUTH_GRPC_PORT` | 50051 | Порт gRPC сервера Auth |
| `AUTH_GRPC_ADDR` | localhost:50051 | Адрес Auth gRPC сервера |
| `TASKS_PORT` | 8082 | Порт HTTP сервера Tasks |
| `TASKS_BASE_URL` | http://193.233.175.221:8082 | Базовый URL Tasks сервиса |
| `HTTPS_GATEWAY` | https://193.233.175.221:8443 | HTTPS эндпоинт через NGINX |
| `DB_HOST` | postgres | Хост PostgreSQL |
| `DB_NAME` | tasksdb | Имя базы данных |
| `DB_PORT` | 5432 | Порт PostgreSQL |
| `DB_USER` | - | Пользователь БД |
| `DB_PASSWORD` | - | Пароль БД |
| `DB_SSLMODE` | disable | Режим SSL для БД |
## Секреты GitHub Actions
| Секрет | Назначение |
|------------|----------------------|
| `SSH_PRIVATE_KEY` | Приватный ключ для доступа к VPS |
| `VPS_HOST` | IP-адрес или домен сервера |
| `VPS_USER` | Пользователь на сервере |
| `DB_PASSWORD` | Пароль для PostgreSQL |
| `TELEGRAM_BOT_TOKEN` | Токен Telegram бота |
| `TELEGRAM_CHAT_ID	ID` | чата для уведомлений |
## Встроенные переменные GitHub
| Переменная | Значение |
|------------|----------------------|
| `${{ github.repository }}` | owner/repo-name |
| `${{ github.actor }}` | Имя пользователя |
| `${{ github.sha }}` | Хэш коммита |
| `${{ github.ref_name }}` | Имя ветки |
| `secrets.GHCR_TOKEN` | Автоматический токен для доступа к API |
## Особенности реализации
- **JSON формат** - все логи выводятся в структурированном JSON виде
- **Прокидывание request-id** - передается в gRPC метаданных при межсервисных вызовах
- **Безопасность** - токены и пароли не попадают в логи

## Уровни логирования
| Уровень | Назначение |
|---------|------------|
| DEBUG | Детальная информация для отладки |
| INFO | Успешные операции |
| WARN | Проблемы, не влияющие на работу |
| ERROR | Критические ошибки |

## Особенности работы с HTTPS
- Самоподписанный сертификат, поэтому в Postman нужно отключить проверку SSL (Settings -> SSL certificate verification: OFF)
- Все запросы идут через NGINX на порт 8443
- NGINX проксирует запросы в соответствующие сервисы (auth и tasks)
- В production-среде используются сертификаты от доверенных центров сертификации

## Ошибки gRPC
| Код | Описание |
|-----|----------|
| 401 Unauthenticated | Невалидный токен |
| 503 DeadlineExceeded | Сервис завис |
| 502 Unavailable | Auth сервис недоступен |
| 500 Internal | Внутренняя ошибка |

## Кэширование (Redis)
### Стратегия cache-aside
1. **GET /v1/tasks/{id}**
   - Проверка кэша по ключу `tasks:task:{id}`
   - Если найден → возврат из кэша
   - Если не найден → запрос в БД → сохранение в кэш → возврат
2. **GET /v1/tasks**
   - Проверка кэша по ключу `tasks:list:{subject}`
   - Если найден → возврат списка из кэша
   - Если не найден → запрос в БД → сохранение в кэш → возврат

### Инвалидация
- **POST /v1/tasks** → удаление `tasks:list:{subject}`
- **PATCH /v1/tasks/{id}** → удаление `tasks:task:{id}` и `tasks:list:{subject}`
- **DELETE /v1/tasks/{id}** → удаление `tasks:task:{id}` и `tasks:list:{subject}`
### TTL с jitter
- Базовый TTL: 120 секунд
- Jitter: случайное значение 0-30 секунд
- Итоговый TTL: 120-150 секунд
### Деградация
При недоступности Redis:
- Ошибки логируются как WARN
- Запросы обслуживаются напрямую из БД
- Сервис продолжает работать

# API Endpoints
## Auth Service (/v1/auth)
### POST http://193.233.175.221:8081/v1/auth/login
- Получение токена доступа
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Body (raw):
```json
{
  "username": "student",
  "password": "student"
}
```
Ответ 200:
```json
{
  "access_token": "demo-token-for-student",
  "token_type": "Bearer"
}
```
Ответ 400:
```json
{
  "error": "invalid request format"
}
```
Ответ 401:
```json
{
  "error": "invalid credentials"
}
```
### GET http://193.233.175.221:8081/v1/auth/verify
- Проверка валидности токена
- Headers:
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Authorization: Bearer Token demo-token-for-student

Ответ 200:
```json
{
  "valid": true,
  "subject": "student"
}
```
Ответ 401:
```json
{
  "valid": false,
  "error": "unauthorized"
}
```
## Tasks Service (/v1/tasks)
### POST 
#### Базовый http://193.233.175.221:8082/v1/tasks
#### С поддержкой HTTPS https://193.233.175.221:8443/v1/tasks
- Создание новой задачи
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Authorization: Bearer Token demo-token-for-student
- Body (raw):
```json
{
  "title": "Do PZ17",
  "description": "split services",
  "due_date": "2026-01-10"
}
```
Ответ 201:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Do PZ17",
  "description": "split services",
  "due_date": "2026-01-10",
  "done": false
}
```
Ошибки:
- 400: Неверный формат запроса
- 401: Неавторизованный запрос (отсутствие или недействительный токен)

### GET
#### Базовый http://193.233.175.221:8082/v1/tasks
#### С поддержкой HTTPS https://193.233.175.221:8443/v1/tasks
- Получение списка всех задач
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Authorization: Bearer Token demo-token-for-student

Ответ 200:
```json
[
    {
        "id": "t_100236.9",
        "title": "Do PZ17",
        "done": false
    },
    {
        "id": "t_100447.8",
        "title": "Do PZ18",
        "done": false
    }
]
```
Ошибки:
- 401: Неавторизованный запрос
### GET (/tasks/search) ДЕМОНСТРАЦИЯ SQL-ИНЪЕКЦИЙ
#### https://193.233.175.221:8443/v1/tasks/search?q={term}&vulnerable=true
- Возвращает задачи, содержащие "{term}" в названии. Нормальный поиск
- Authorization: Bearer Token demo-token-for-student
#### SQL-инъекция https://193.233.175.221:8443/v1/tasks/search?q=' OR '1'='1&vulnerable=true
- Возвращает ВСЕ задачи из базы данных, игнорируя фильтр по пользователю
#### Деструктивная SQL-инъекция GET https://193.233.175.221:8443/v1/tasks/search?q='; DROP TABLE tasks; --&vulnerable=true
- Удаляет таблицу
#### Безопасная SQL-инъекция https://193.233.175.221:8443/v1/tasks/search?q={term}
- Ищет задачи, содержащие буквально строку "' OR '1'='1" в названии
- Authorization: Bearer Token demo-token-for-student
```json
[
    {
        "id": "t_100236.9",
        "title": "Do PZ17",
        "done": false
    },
    {
        "id": "t_100447.8",
        "title": "Do PZ18",
        "done": false
    }
]
```
### GET http://193.233.175.221:8082/v1/tasks/{id}
- Получение задачи по ID
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Authorization: Bearer Token demo-token-for-student

Ответ 200:
```json
{
    "id": "t20260219102056",
    "title": "Do PZ17",
    "description": "split services",
    "due_date": "2026-01-10",
    "done": false
}
```
Ответ 404:
```json
{
  "error": "task not found"
}
```
Ошибки:
- 401: Неавторизованный запрос

### PATCH http://193.233.175.221:8082/v1/tasks/{id}
- Обновление задачи
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Authorization: Bearer Token demo-token-for-student
- Body (raw):
```json
{
  "title": "Do PZ17 (updated)",
  "done": true
}
```
Ответ 200:
```json
{
    "id": "t20260219102056",
    "title": "Do PZ17 (updated)",
    "description": "split services",
    "due_date": "2026-01-10",
    "done": true
}
```
Ошибки:
- 400: Неверный формат запроса
- 404: Задача не найдена
- 401: Неавторизованный запрос

### DELETE http://193.233.175.221:8082/v1/tasks/{id}
- Удаление задачи
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Authorization: Bearer Token demo-token-for-student

Ответ:
- 204: Успешное удаление, тела ответа нет

Ошибки:
- 404: Задача не найдена
- 401: Неавторизованный запрос

### Общие коды ошибок для Tasks Service
- 400 Bad Request           неверный формат запроса
- 401 Unauthorized          отсутствует или недействительный токен
- 404 Not Found             задача не найдена
- 500 Internal Server Error внутренняя ошибка сервиса
- 502 Bad Gateway           недоступен Auth сервис
- 503 Service Unavailable   таймаут при обращении к Auth

## Task Service HTTPS
### GET (/health)
Проверка состояния работы сервиса

Ответ 200:
```json
{
  "status": "ok",
  "service": "tasks"
}
```
### GET (/metrics)
```
HTTPS Gateway is working!
```
## CSRF защита
### POST https://193.233.175.221:8443/v1/auth/login
- Логин и получение cookies
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Body (raw):
```json
{
  "username": "student",
  "password": "student"
}
```
Ответ 200:
```json
{
  "message": "Login successful"
}
```
### POST https://193.233.175.221:8443/v1/tasks
- Создание задачи
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
    - X-CSRF-Token: [значение из cookie csrf_token]
    - Authorization: Bearer Token demo-token-for-student
- Cookies:
    - session_id (из предыдущего логина)
    - csrf_token (из предыдущего логина)
- Body (raw):
```json
{
  "title": "CSRF Safe",
  "description": "Should succeed",
  "due_date": "2026-03-01"
}
```
Ответ 201:
```json
{
  "id": "t20260227123456",
  "title": "CSRF Safe",
  "description": "Should succeed",
  "due_date": "2026-03-01",
  "done": false
}
```
Ответ 403:
```json
{
  "error": "CSRF token missing in header"
}
```
Ошибки:
- 403: CSRF token missing/invalid - отсутствует или неверный CSRF токен
- 401: Unauthorized - отсутствует сессия

### GET https://193.233.175.221:8443/v1/tasks
- Создание задачи
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
    - Authorization: Bearer Token demo-token-for-student
- Cookies:
    - session_id (из предыдущего логина)
```
Ответ 200:
```json
[
  {
    "id": "t20260227123456",
    "title": "CSRF Safe",
    "done": false
  }
]
```
### POST https://193.233.175.221:8443/v1/auth/logout
- Выход из системы (очистка cookies)
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
- Cookies:
    - session_id (из предыдущего логина)
```
Ответ 200:
```json
{
  "message": "Logout successful"
}
```
### GET https://193.233.175.221:8443/v1/auth/csrf
- Получение CSRF токена (если нужно обновить)
- Headers:
    - Content-Type: application/json
    - X-Request-ID: test-123 (опционально, но рекомендуется)
    - Authorization: Bearer Token demo-token-for-student
- Cookies:
    - session_id (из предыдущего логина)
```
Ответ 200:
```json
{
  "csrf_token": "base64-encoded-csrf-token"
}
```
Ошибки:
- 401: Unauthorized - отсутствует или невалидная сессия

## GraphQL API
### POST http://193.233.175.221:8090/query
- Headers:
    - Content-Type: application/json


### Playground http://193.233.175.221:8090/
#### Task
| Поле | Тип | Описание |
|------|-----|----------|
| id | ID! | Уникальный идентификатор |
| title | String! | Название задачи |
| description | String | Описание задачи |
| due_date | String | Дата выполнения |
| done | Boolean! | Статус выполнения |
| created_at | String | Дата создания |
| updated_at | String | Дата обновления |

#### CreateTaskInput
| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| title | String! | Да | Название задачи |
| description | String | Нет | Описание задачи |
| due_date | String | Нет | Дата выполнения |

#### UpdateTaskInput
| Поле | Тип | Описание |
|------|-----|----------|
| title | String | Новое название |
| description | String | Новое описание |
| due_date | String | Новая дата |
| done | Boolean | Новый статус |

### Запросы (Queries)
#### Получить все задачи
```graphql
query {
  tasks {
    id
    title
    description
    due_date
    done
    created_at
    updated_at
  }
}
```
Ответ:
```json
{
  "data": {
    "tasks": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "GraphQL Task",
        "description": "Created via GraphQL",
        "due_date": "2026-03-15",
        "done": false,
        "created_at": "2026-03-08T12:00:00Z",
        "updated_at": "2026-03-08T12:00:00Z"
      }
    ]
  }
}
```
#### Получить задачу по ID
```graphql
query GetTask($id: ID!) {
  task(id: $id) {
    id
    title
    description
    done
  }
}
```
Variables:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```
#### Создать задачу
```graphql
mutation CreateTask($input: CreateTaskInput!) {
  createTask(input: $input) {
    id
    title
    description
    done
    created_at
  }
}
```
Variables:
```json
{
  "input": {
    "title": "New Task",
    "description": "Task description",
    "due_date": "2026-03-15"
  }
}
```
Ответ:
```json
{
  "data": {
    "createTask": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "New Task",
      "description": "Task description",
      "done": false,
      "created_at": "2026-03-08T12:00:00Z"
    }
  }
}
```
#### Обновить задачу
```graphql
mutation UpdateTask($id: ID!, $input: UpdateTaskInput!) {
  updateTask(id: $id, input: $input) {
    id
    title
    done
    updated_at
  }
}
```
Variables:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "input": {
    "done": true,
    "title": "Updated Task"
  }
}
```
#### Удалить задачу
```graphql
mutation DeleteTask($id: ID!) {
  deleteTask(id: $id)
}
```
Variables:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```
Ответ:
```json
{
  "data": {
    "deleteTask": true
  }
}
```
Ошибки:
- 200: Успешный запрос
- 400: Неверный синтаксис GraphQL
- 500: Внутренняя ошибка сервера