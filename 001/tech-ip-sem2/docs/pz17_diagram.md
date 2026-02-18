# Диаграмма взаимодействия микросервисов

## Схема последовательности
# Диаграмма взаимодействия микросервисов

## Схема последовательности

```mermaid
sequenceDiagram
    participant C as Client
    participant T as Tasks Service (:8082)
    participant A as Auth Service (:8081)
    
    Note over C,A: 1. Получение токена
    C->>A: POST /v1/auth/login
    Note right of A: Проверка username/password
    A-->>C: 200 OK (access_token)
    
    Note over C,T: 2. Работа с задачами
    C->>T: POST /v1/tasks + Authorization
    Note right of T: Извлечение токена
    
    T->>A: GET /v1/auth/verify + Token
    Note right of A: Валидация токена
    A-->>T: 200 OK (valid)
    
    Note right of T: Создание задачи
    T-->>C: 201 Created (task data)
    
    Note over C,T: 3. Ошибка авторизации
    C->>T: GET /v1/tasks + Invalid Token
    T->>A: GET /v1/auth/verify
    A-->>T: 401 Unauthorized
    T-->>C: 401 Unauthorized