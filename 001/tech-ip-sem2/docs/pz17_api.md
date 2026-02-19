# Примеры запросов

## API Endpoints

### Auth Service (/v1/auth)

#### POST http://193.233.175.221:8081/v1/auth/login
- Получение токена доступа
- Запрос:
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
#### GET http://193.233.175.221:8081/v1/auth/verify
- Проверка валидности токена
- Заголовки:
Authorization: Bearer <token>
X-Request-ID: <uuid> (опционально)

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
### Tasks Service (/v1/tasks)

#### POST http://193.233.175.221:8081/v1/tasks
- Создание новой задачи
- Заголовки:
Authorization: Bearer <token>
X-Request-ID: <uuid> (рекомендуется)
- Запрос:
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
  "id": "t_001",
  "title": "Do PZ17",
  "description": "split services",
  "due_date": "2026-01-10",
  "done": false
}
```
Ошибки:
- 400: Неверный формат запроса
- 401: Неавторизованный запрос (отсутствие или недействительный токен)

#### GET http://193.233.175.221:8081/v1/tasks
- Получение списка всех задач
- Заголовки:
Authorization: Bearer <token>
X-Request-ID: <uuid> (рекомендуется)

Ответ 200:
```json
[
  {
    "id": "t_001",
    "title": "Do PZ17",
    "description": "split services",
    "due_date": "2026-01-10",
    "done": false
  },
  {
    "id": "t_002",
    "title": "Read lecture",
    "description": "",
    "due_date": "2025-12-01",
    "done": true
  }
]
```
Ошибки:
- 401: Неавторизованный запрос

#### GET http://193.233.175.221:8081/v1/tasks/{id}
- Получение задачи по ID
- Заголовки:
Authorization: Bearer <token>
X-Request-ID: <uuid> (рекомендуется)

Ответ 200:
```json
{
  "id": "t_001",
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

#### PATCH http://193.233.175.221:8081/v1/tasks/{id}
- Обновление задачи
- Заголовки:
Authorization: Bearer <token>
X-Request-ID: <uuid> (рекомендуется)
-Запрос:
```json
{
  "title": "Do PZ17 (updated)",
  "done": true
}
```
Ответ 200:
```json
{
  "id": "t_001",
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

#### DELETE http://193.233.175.221:8081/v1/tasks/{id}
- Удаление задачи
- Заголовки:
Authorization: Bearer <token>
X-Request-ID: <uuid> (рекомендуется)

Ответ:
- 204: Успешное удаление, тела ответа нет

Ошибки:
- 404: Задача не найдена
- 401: Неавторизованный запрос

Общие коды ошибок для Tasks Service:
- 400 Bad Request неверный формат запроса
- 401 Unauthorized отсутствует или недействительный токен
- 404 Not Found задача не найдена
- 500 Internal Server Error внутренняя ошибка сервиса
- 502 Bad Gateway недоступен Auth сервис
- 503 Service Unavailable таймаут при обращении к Auth