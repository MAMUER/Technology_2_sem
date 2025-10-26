# Практическая работа №10
# Николаенко Михаил ЭФМО-02-21

## Описание проекта и требования

JWT Authentication Microservice - это REST API сервис аутентификации и авторизации на Go. Сервис реализует stateless JWT-аутентификацию с access/refresh токенами, RBAC авторизацией и ABAC правилами доступа.
Основные возможности: безопасный login/logout, обновление токенов, управление правами доступа на основе ролей (admin/user) и атрибутов пользователя. Сервис использует chi роутер, middleware-цепочки для аутентификации и современные практики безопасности включая bcrypt для хранения паролей.

Для работы с командой make в PowerShell необходимо установить менеджер пакетов Chocolatey и установить команду make

Проект на языке Go (необходима версия 1.21 и выше) с REST-API

- `POST /api/v1/login` - Авторизация пользователя
- `POST /api/v1/refresh` - Токены обновления
- `POST /api/v1/logout` - Выход пользователя
- `GET  /api/v1/me` - Получить текущего пользователя
- `GET  /api/v1/users/{id}` - Получить пользователя ID (ABAC защита)
- `GET  /api/v1/admin/stats` - Admin статистика (RBAC защита)"

## Команды запуска/сборки

### Сборка приложения:

make build

### Запуск приложения:

make run

### Инструкция:

make help

## Команды:

### Логин — получить токен админа
curl -s -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d "{\"Email\":\"admin@example.com\",\"Password\":\"secret123\"}"

Ответ:

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwejEwLWNsaWVudHMiLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwiZXhwIjoxNzYxMjE2MzY1LCJpYXQiOjE3NjEyMDkxNjUsImlzcyI6InB6MTAtYXV0aCIsInJvbGUiOiJhZG1pbiIsInN1YiI6MX0.GqjQ13GOvySLMs1CIcst7Qf2jBnH-EXc8euAEGDnGJ8","user":{"email":"admin@example.com","id":1,"role":"admin"}}

### Сохранение токена в переменную для удобства (работает не для всех консолей)
set TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwejEwLWNsaWVudHMiLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwiZXhwIjoxNzYxNDAyNTYzLCJpYXQiOjE3NjE0MDE2NjMsImlzcyI6InB6MTAtYXV0aCIsInJvbGUiOiJhZG1pbiIsInN1YiI6MSwidHlwZSI6ImFjY2VzcyJ9.V8rPOcpM9-LJIn4zMijt5tVF8RjX0Sa_ju7gYr4YtPI

### Доступ к защищённым ручкам:
curl -s http://localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN"

curl -s http://localhost:8080/api/v1/admin/stats -H "Authorization: Bearer $TOKEN"

Ответы:

{"email":"admin@example.com","id":1,"role":"admin"}

{"stats":"admin only data","user":{"email":"admin@example.com","id":1,"role":"admin"},"users":42}

### Логин — получить токен пользователя

curl -s -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d "{\"Email\":\"user@example.com\",\"Password\":\"secret123\"}"

Ответ:

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwejEwLWNsaWVudHMiLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NjEyMTg0MDIsImlhdCI6MTc2MTIxMTIwMiwiaXNzIjoicHoxMC1hdXRoIiwicm9sZSI6InVzZXIiLCJzdWIiOjJ9.i_RDZ-PhsO1JthNOS7uR4HweUXZ_YYzO-cEAKc7SKqE","user":{"email":"user@example.com","id":2,"role":"user"}}

### Доступ к защищённым ручкам:
curl -i http://localhost:8080/api/v1/admin/stats -H "Authorization: Bearer $TOKEN_USER"

Ответ:

HTTP/1.1 403 Forbidden
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 25 Oct 2025 15:16:16 GMT
Content-Length: 38

{"error": "insufficient permissions"}

## Тесты (Для PowerShell)
### 1. Логин админа
$admin = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/login" -Method POST -ContentType "application/json" -Body '{"email":"admin@example.com","password":"secret123"}'
$ADMIN_ACCESS = $admin.access_token
$ADMIN_REFRESH = $admin.refresh_token
Write-Host "Admin Access: $ADMIN_ACCESS"
Write-Host "Admin Refresh: $ADMIN_REFRESH"

### 2. Логин пользователя
$user = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/login" -Method POST -ContentType "application/json" -Body '{"email":"user@example.com","password":"secret123"}'
$USER_ACCESS = $user.access_token
$USER_REFRESH = $user.refresh_token
Write-Host "User Access: $USER_ACCESS"
Write-Host "User Refresh: $USER_REFRESH"

### 3. Тест /me для админа
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/me" -Headers @{"Authorization"="Bearer $ADMIN_ACCESS"}

### 4. Тест /me для пользователя
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/me" -Headers @{"Authorization"="Bearer $USER_ACCESS"}

### 5. ABAC тест: пользователь запрашивает свой профиль (должен работать)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/2" -Headers @{"Authorization"="Bearer $USER_ACCESS"}

### 6. ABAC тест: пользователь запрашивает чужой профиль (должен вернуть 403)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/1" -Headers @{"Authorization"="Bearer $USER_ACCESS"}

### 7. ABAC тест: админ запрашивает любой профиль (должен работать)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/2" -Headers @{"Authorization"="Bearer $ADMIN_ACCESS"}

### 8. Тест админского эндпоинта
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/stats" -Headers @{"Authorization"="Bearer $ADMIN_ACCESS"}

### 9. Тест refresh токена
$body = @{refresh_token = $USER_REFRESH} | ConvertTo-Json
$refresh = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/refresh" -Method POST -ContentType "application/json" -Body $body
$NEW_ACCESS = $refresh.access_token
Write-Host "New Access: $($NEW_ACCESS)"

### 10. Тест нового access токена
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/me" -Headers @{"Authorization"="Bearer $NEW_ACCESS"}

### 11. Логаут
$logoutBody = @{refresh_token = $USER_REFRESH} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/logout" -Method POST -ContentType "application/json" -Body $logoutBody
Write-Host "Logout successful"

## Структура проекта
```
C:.
└───pz10-auth
    │   .env
    │   go.mod
    │   go.sum
    │   Makefile
    │   README.md
    │
    ├───bin
    │       server.exe
    │
    ├───cmd
    │   └───server
    │           main.go
    │
    ├───internal
    │   ├───core
    │   │       service.go
    │   │       user.go
    │   │
    │   ├───http
    │   │   │   router.go
    │   │   │
    │   │   └───middleware
    │   │           authn.go
    │   │           authz.go
    │   │
    │   ├───platform
    │   │   ├───config
    │   │   │       config.go
    │   │   │
    │   │   └───jwt
    │   │           jwt.go
    │   │
    │   └───repo
    │           refresh_mem.go
    │           user_mem.go
    │
    └───PR10
```
## Переменные окружения (.env)

APP_PORT=8080 - Порт сервера

JWT_SECRET=your-secret-key - JWT пароль (необходим)

JWT_ACCESS_TTL=15m - Токен доступа TTL

JWT_REFRESH_TTL=168h - Токен обновления TTL (7 days)

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR10/Screenshot_1.png)

Проверка и запуск приложения (+ логи)

![фото2](./PR10/Screenshot_2.png)

![фото3](./PR10/Screenshot_3.png)

Регистрация админа и пользователя, затем входы в систему и проверка доступа фукнций

![фото4](./PR10/Screenshot_4.png)

Проверка новых добавленных функций:

![фото5](./PR10/Screenshot_5.png)

![фото6](./PR10/Screenshot_6.png)

![фото7](./PR10/Screenshot_7.png)

![фото8](./PR10/Screenshot_8.png)

![фото9](./PR10/Screenshot_9.png)

![фото10](./PR10/Screenshot_10.png)

![фото11](./PR10/Screenshot_11.png)

![фото12](./PR10/Screenshot_12.png)

Структура проекта

![фото13](./PR10/Screenshot_13.png)