3. Спецификация API
3.1. Auth Service (порт 8081)
POST /v1/auth/login
Получение токена доступа

Request:

json
{
    "username": "student",
    "password": "student"
}
Response 200:

json
{
    "access_token": "demo-token-for-student",
    "token_type": "Bearer"
}
Response 400:

json
{
    "error": "invalid request format"
}
Response 401:

json
{
    "error": "invalid credentials"
}
GET /v1/auth/verify
Проверка валидности токена

Headers:

Authorization: Bearer <token>

X-Request-ID: <uuid> (опционально)

Response 200:

json
{
    "valid": true,
    "subject": "student"
}
Response 401:

json
{
    "valid": false,
    "error": "unauthorized"
}
3.2. Tasks Service (порт 8082)
Все запросы требуют заголовок:

Authorization: Bearer <token>

X-Request-ID: <uuid> (рекомендуется)

POST /v1/tasks
Создание новой задачи

Request:

json
{
    "title": "Do PZ17",
    "description": "split services",
    "due_date": "2026-01-10"
}
Response 201:

json
{
    "id": "t_001",
    "title": "Do PZ17",
    "description": "split services",
    "due_date": "2026-01-10",
    "done": false
}
GET /v1/tasks
Получение списка всех задач

Response 200:

json
[
    {
        "id": "t_001",
        "title": "Do PZ17",
        "done": false
    },
    {
        "id": "t_002",
        "title": "Read lecture",
        "done": true
    }
]
GET /v1/tasks/{id}
Получение задачи по ID

Response 200:

json
{
    "id": "t_001",
    "title": "Do PZ17",
    "description": "split services",
    "due_date": "2026-01-10",
    "done": false
}
Response 404:

json
{
    "error": "task not found"
}
PATCH /v1/tasks/{id}
Обновление задачи

Request:

json
{
    "title": "Do PZ17 (updated)",
    "done": true
}
Response 200:

json
{
    "id": "t_001",
    "title": "Do PZ17 (updated)",
    "description": "split services",
    "due_date": "2026-01-10",
    "done": true
}
DELETE /v1/tasks/{id}
Удаление задачи

Response 204 (без тела)

Коды ошибок для Tasks Service:
400 Bad Request - неверный формат запроса

401 Unauthorized - токен невалиден или отсутствует

404 Not Found - задача не найдена

500 Internal Server Error - внутренняя ошибка сервиса

502 Bad Gateway - Auth service недоступен

503 Service Unavailable - таймаут при обращении к Auth