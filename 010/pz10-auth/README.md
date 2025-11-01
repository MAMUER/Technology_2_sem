# Практическая работа №10
# Николаенко Михаил ЭФМО-02-21

## Описание проекта

JWT Authentication Microservice - это REST API сервис аутентификации и авторизации на Go. Сервис реализует stateless JWT-аутентификацию с access/refresh токенами, RBAC авторизацией и ABAC правилами доступа.
Основные возможности: безопасный login/logout, обновление токенов, управление правами доступа на основе ролей (admin/user) и атрибутов пользователя. Сервис использует chi роутер, middleware-цепочки для аутентификации и современные практики безопасности включая bcrypt для хранения паролей.

## Требования
- Go версии 1.25 и выше

## Основные эндпоинты

### Аутентификация
#### Авторизация админа
- `POST http://193.233.175.221:8084/api/v1/login`
  - `Headers` Key: Content-Type Value: application/json
  - `Body`: {"Email": "admin@example.com","Password": "secret123"}

#### Авторизация пользователя  
- `POST http://193.233.175.221:8084/api/v1/login`
  - `Headers` Key: Content-Type Value: application/json
  - `Body`: {"Email": "user@example.com","Password": "secret123"}

#### Авторизация пользователя 2
- `POST http://193.233.175.221:8084/api/v1/login`
  - `Headers` Key: Content-Type Value: application/json
  - `Body`: {"Email": "user2@example.com","Password": "secret123"}

#### Обновление токена
- `POST http://193.233.175.221:8084/api/v1/refresh`
  - `Headers` Key: Content-Type Value: application/json
  - `Body`: {"refresh_token": "token{}"}

#### Выход пользователя
- `POST http://193.233.175.221:8084/api/v1/logout`
  - `Headers` Key: Content-Type Value: application/json
  - `Body`: {"refresh_token": "token"}

### Пользовательские эндпоинты (требуют Authorization header)
#### Получить текущего пользователя
- `GET http://193.233.175.221:8084/api/v1/me`
  - `Authorization` `Bearer Token` {token}

#### Получить пользователя по ID (ABAC защита)
- `GET http://193.233.175.221:8084/api/v1/users/{id}`
  - `Authorization` `Bearer Token` {token}

### Админские эндпоинты (RBAC защита)
#### Статистика системы (только для админа)
- `GET http://193.233.175.221:8084/api/v1/admin/stats`
  - `Authorization` `Bearer Token` {token}

#### Получить любого пользователя (только для админа)
- `GET http://193.233.175.221:8084/api/v1/users/{id}`
  - `Authorization` `Bearer Token` {token}

## Команды

### Логин — получить токен админа
http://193.233.175.221:8084/api/v1/login

Ответ:

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwejEwLWNsaWVudHMiLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwiZXhwIjoxNzYxMjE2MzY1LCJpYXQiOjE3NjEyMDkxNjUsImlzcyI6InB6MTAtYXV0aCIsInJvbGUiOiJhZG1pbiIsInN1YiI6MX0.GqjQ13GOvySLMs1CIcst7Qf2jBnH-EXc8euAEGDnGJ8","user":{"email":"admin@example.com","id":1,"role":"admin"}}

### Доступ к защищённым ручкам:
http://193.233.175.221:8084/api/v1/me

http://193.233.175.221:8084/api/v1/admin/stats

Ответы:

{"email":"admin@example.com","id":1,"role":"admin"}

{"stats":"admin only data","user":{"email":"admin@example.com","id":1,"role":"admin"},"users":42}

### Логин — получить токен пользователя

http://193.233.175.221:8084/api/v1/login

Ответ:

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwejEwLWNsaWVudHMiLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NjEyMTg0MDIsImlhdCI6MTc2MTIxMTIwMiwiaXNzIjoicHoxMC1hdXRoIiwicm9sZSI6InVzZXIiLCJzdWIiOjJ9.i_RDZ-PhsO1JthNOS7uR4HweUXZ_YYzO-cEAKc7SKqE","user":{"email":"user@example.com","id":2,"role":"user"}}

### Доступ к защищённым ручкам:
http://193.233.175.221:8084/api/v1/admin/stats

Ответ:

HTTP/1.1 403 Forbidden
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 25 Oct 2025 15:16:16 GMT
Content-Length: 38

{"error": "insufficient permissions"}

## Тесты (Для PowerShell)
### 1. Логин админа
$admin = Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/login" -Method POST -ContentType "application/json" -Body '{"email":"admin@example.com","password":"secret123"}'
$ADMIN_ACCESS = $admin.access_token
$ADMIN_REFRESH = $admin.refresh_token
Write-Host "Admin Access: $ADMIN_ACCESS"
Write-Host "Admin Refresh: $ADMIN_REFRESH"

### 2. Логин пользователя
$user = Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/login" -Method POST -ContentType "application/json" -Body '{"email":"user@example.com","password":"secret123"}'
$USER_ACCESS = $user.access_token
$USER_REFRESH = $user.refresh_token
Write-Host "User Access: $USER_ACCESS"
Write-Host "User Refresh: $USER_REFRESH"

### 3. Тест /me для админа
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/me" -Headers @{"Authorization"="Bearer $ADMIN_ACCESS"}

### 4. Тест /me для пользователя
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/me" -Headers @{"Authorization"="Bearer $USER_ACCESS"}

### 5. ABAC тест: пользователь запрашивает свой профиль (должен работать)
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/users/2" -Headers @{"Authorization"="Bearer $USER_ACCESS"}

### 6. ABAC тест: пользователь запрашивает чужой профиль (должен вернуть 403)
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/users/1" -Headers @{"Authorization"="Bearer $USER_ACCESS"}

### 7. ABAC тест: админ запрашивает любой профиль (должен работать)
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/users/2" -Headers @{"Authorization"="Bearer $ADMIN_ACCESS"}

### 8. Тест админского эндпоинта
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/admin/stats" -Headers @{"Authorization"="Bearer $ADMIN_ACCESS"}

### 9. Тест refresh токена
$body = @{refresh_token = $USER_REFRESH} | ConvertTo-Json
$refresh = Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/refresh" -Method POST -ContentType "application/json" -Body $body
$NEW_ACCESS = $refresh.access_token
Write-Host "New Access: $($NEW_ACCESS)"

### 10. Тест нового access токена
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/me" -Headers @{"Authorization"="Bearer $NEW_ACCESS"}

### 11. Логаут
$logoutBody = @{refresh_token = $USER_REFRESH} | ConvertTo-Json
Invoke-RestMethod -Uri "http://193.233.175.221:8084/api/v1/logout" -Method POST -ContentType "application/json" -Body $logoutBody
Write-Host "Logout successful"

## Структура проекта
```
.
├── bin
│   └── server.exe
├── cmd
│   └── server
│       └── main.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── core
│   │   ├── service.go
│   │   └── user.go
│   ├── http
│   │   ├── middleware
│   │   │   ├── authn.go
│   │   │   ├── authz.go
│   │   │   └── contex.go
│   │   └── router.go
│   ├── platform
│   │   ├── config
│   │   │   └── config.go
│   │   └── jwt
│   │       └── jwt.go
│   └── repo
│       ├── refresh_mem.go
│       └── user_mem.go
├── Makefile
├── PR10
└── README.md
```

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR10/Screenshot_1.png)

Проверка и запуск приложения (+ логи)

![фото2](./PR10/Screenshot_2.png)

![фото3](./PR10/Screenshot_3.png)

Регистрация админа и пользователя, затем входы в систему и проверка доступа фукнций, проверка новых добавленных функций:

![фото4](./PR10/Screenshot_4.png)

![фото5](./PR10/Screenshot_5.png)

![фото6](./PR10/Screenshot_6.png)

![фото7](./PR10/Screenshot_7.png)

![фото8](./PR10/Screenshot_8.png)

![фото9](./PR10/Screenshot_9.png)

![фото10](./PR10/Screenshot_10.png)

![фото11](./PR10/Screenshot_11.png)

Структура проекта

![фото13](./PR10/Screenshot_13.png)