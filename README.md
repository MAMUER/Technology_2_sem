# Практическое занятие №1-16
# Группа: ЭФМО-02-25
# Студент: Николаенко Михаил

## Содержание
- [Описание проекта](#описание-проекта)
- [Технологии](#технологии)
- [Архитектура](#архитектура)
- [Документация по практическим занятиям](#документация-по-практическим-занятиям)
- [Команды запуска и сборки](#команды-makefile)
- [Переменные окружения](#переменные-окружения)
- [Сборка и запуск](#сборка-и-запуск)
- [CI/CD Pipeline](#cicd-pipeline)
- [Структура проекта](#структура-проекта)
- [Скриншоты работы](#скриншоты-работы)

## Описание проекта
Проект представляет собой микросервисную архитектуру с двумя сервисами, демонстрирующую современные практики разработки на Go. Реализованы:
- **ПЗ №1-3**: Разделение монолита на микросервисы, логирование (zap), request-id трассировка
- **ПЗ №4**: Метрики, Prometheus, Grafana
- **ПЗ №5**: HTTPS/TLS, PostgreSQL, защита от SQL-инъекций
- **ПЗ №6**: CSRF/XSS защита, безопасные cookies
- **ПЗ №7**: Написание Dockerfile и сборка контейнера
- **ПЗ №8**: Настройка GitHub Actions / GitLab CI для деплоя приложения
- **ПЗ №9**: Реализация распределённого кэша (Redis cluster)
- **ПЗ №10**: Горизонтальное масштабирование: использование Load Balancer (NGINX)
- **ПЗ №11**: Создание GraphQL API с использованием gqlgen. Запросы и мутации
- **ПЗ №12**: Сравнение REST и GraphQL: разработка одного и того же функционала двумя способами
- **ПЗ №13**: Подключение к RabbitMQ. Отправка и получение сообщений
- **ПЗ №14**: Реализация очереди задач (producer–consumer): retries, DLQ, идемпотентность
- **ПЗ №15**: Деплой приложения на VPS. Настройка system
- **ПЗ №16**: Публикация приложения в Kubernetes (минимальный манифест)

### Технологии

| Компонент | Технология | Назначение |
| :--- | :--- | :--- |
| **Язык** | Go 1.25+ | Основной язык разработки |
| **Логгер** | Uber Zap | Структурированное логирование |
| **Метрики** | Prometheus + Grafana | Сбор и визуализация метрик |
| **База данных** | PostgreSQL | Хранение данных (с параметризацией) |
| **Кэширование** | Redis | Распределённый кэш (cache-aside pattern) |
| **Очереди сообщений** | RabbitMQ | Асинхронная обработка событий и задач |
| **gRPC** | Protocol Buffers | Межсервисное взаимодействие |
| **REST API** | net/http + middleware | CRUD операции, авторизация, логирование |
| **GraphQL API** | gqlgen | Альтернативный API с гибкой выборкой полей |
| **TLS/HTTPS** | NGINX | Терминирование HTTPS, reverse proxy |
| **Балансировка** | NGINX | Распределение трафика между репликами |
| **Контейнеризация** | Docker + Docker Compose | Запуск сервисов, multi-stage сборка |
| **Оркестрация** | Kubernetes (minikube) | Деплой и управление контейнерами |
| **Безопасность** | CSRF Double Submit, XSS sanitization | Защита от веб-уязвимостей |
| **Трассировка** | Request ID middleware | Сквозная трассировка запросов |
| **Идемпотентность** | In-memory storage | Защита от дублирования сообщений |
| **Dead Letter Queue** | RabbitMQ DLX/DLQ | Хранение необработанных сообщений |
| **Retry механизм** | RabbitMQ TTL + retry queue | Повторные попытки обработки с задержкой |
| **CI/CD** | GitHub Actions | Автоматизация тестирования, сборки и деплоя |
| **Реестр образов** | GitHub Container Registry (GHCR) | Хранение Docker образов |
| **Управление процессами** | systemd | Управление сервисом на VPS |
| **Логирование инфраструктуры** | journalctl | Просмотр логов systemd сервисов |

## Архитектура
### Auth Service (порт 8081 HTTP, 50051 gRPC)
- Аутентификация и выдача токенов
- Валидация токенов для других сервисов
- Управление сессиями и secure cookies
- CSRF защита (Double Submit Cookie)
- Хранение пользователей в памяти

### Tasks Service (порт 8082 HTTP)
- CRUD операции с задачами
- Хранение в памяти или PostgreSQL
- Проверка доступа через Auth Service (gRPC)
- Request-id трассировка
- CSRF защита для опасных методов
- XSS защита (санитизация ввода)
- Кэширование через Redis (cache-aside pattern)
- Инвалидация кэша при изменении данных

### GraphQL Service (порт 8090 HTTP)
- GraphQL API для задач
- Playground для тестирования запросов
- Поддержка Query и Mutation
- Интеграция с PostgreSQL

### Worker Service
- 2 экземпляра (worker-1, worker-2)
- Потребление событий из RabbitMQ
- Подтверждение обработки (ack)
- Prefetch = 1 для контроля нагрузки

### Система очередей задач (Job Queue)
| Очередь | Назначение | Особенности |
| :--- | :--- | :--- |
| `task_jobs` | Основная очередь | Обработка задач |
| `task_jobs_retry` | Повторные попытки | TTL 10 сек, возврат в основную |
| `task_jobs_dlq` | Dead Letter Queue | Необработанные сообщения |

- **Максимум попыток:** 3
- **Идемпотентность:** хранение обработанных `message_id` в памяти

### Мониторинг
| Компонент | Порт | Назначение |
| :--- | :--- | :--- |
| **Prometheus** | 9090 | Сбор метрик |
| **Grafana** | 3000 | Визуализация метрик |

- **Метрики:** RPS, ошибки, длительность запросов, активные запросы

### Базы данных и очереди
| Компонент | Порт | Назначение |
| :--- | :--- | :--- |
| **PostgreSQL** | 5432 | Основное хранилище данных |
| **Redis** | 6379 | Кэширование |
| **RabbitMQ** | 5672 (AMQP), 15672 (UI) | Очереди сообщений |

### HTTPS Gateway и Балансировка
| Компонент | Порт | Назначение |
| :--- | :--- | :--- |
| **NGINX (Gateway)** | 8443 (HTTPS) | Терминирование SSL/TLS, проксирование |
| **NGINX (LB)** | 8080 | Балансировка трафика (3 реплики tasks) |

- **Балансировка:** Round-robin распределение
- **Health checks:** через `/health`
- **Сертификаты:** Самоподписанные для разработки

### CI/CD Pipeline (GitHub Actions)
- **Test:** линтинг и тестирование кода
- **Build:** сборка бинарников для всех сервисов
- **Docker:** сборка и публикация образов в GHCR
- **Deploy:** деплой на VPS через Docker Compose
- **Notify:** уведомления в Telegram

### Инфраструктура развёртывания

#### VPS
- **systemd:** управление сервисом tasks (порт 9082)
- **journalctl:** централизованное логирование
- **Автозапуск:** при старте системы
- **Self-healing:** автоматический перезапуск при падении

#### Kubernetes (локальный стенд)
- **Minikube:** локальный кластер
- **Deployment:** 2-3 реплики tasks
- **Service:** ClusterIP для внутреннего доступа
- **ConfigMap:** конфигурация сервиса
- **Secret:** чувствительные данные
- **Probes:** readiness и liveness проверки
- **Port-forward:** доступ из локальной сети

## Документация по практическим занятиям
### Практические занятия №1-3 (Логирование)
- [**API Endpoints**](docs/pz_api.md) - Полное описание всех API методов
- [**Диаграмма архитектуры**](docs/pz17_diagram.md) - Схема взаимодействия сервисов

**Основные моменты:**
- Выбран логгер **Uber Zap** (JSON формат, производительность)
- Реализован request-id для трассировки через оба сервиса
- Стандарт полей: level, ts, service, request_id, method, path, status, duration_ms

### Практическое занятие №4 (Метрики)
- [**Метрики и мониторинг**](docs/pz20_metrics.md) - Настройка Prometheus/Grafana

**Метрики:**
- `http_requests_total` - счётчик запросов (method, route, status)
- `http_request_duration_seconds` - гистограмма длительности
- `http_requests_in_flight` - текущие активные запросы

**Реализовано:**
- Самоподписанные SSL-сертификаты
- NGINX как TLS-терминатор
- PostgreSQL с параметризованными запросами
- Защита от SQL-инъекций
- Уязвимая версия для демонстрации

### Практическое занятие №6 (CSRF/XSS и безопасные cookies)

**CSRF защита (Double Submit Cookie):**
| Cookie | HttpOnly | Secure | SameSite | Назначение |
|--------|----------|--------|----------|------------|
| session_id | Да | Да | Lax | Идентификатор сессии |
| csrf_token | Нет | Да | Lax | Токен для CSRF |

**XSS защита:**
- Санитизация ввода (`SanitizeText`, `SanitizeHTML`)
- Заголовки безопасности (CSP, X-Frame-Options, HSTS)

### Практическое занятие №7 (Dockerfile и контейнеризация)
#### Пояснение стадий сборки (multi-stage)
**Stage 1: Builder**
| Шаг | Назначение |
|-----|------------|
| `FROM golang:1.25-alpine` | Базовый образ с Go для компиляции |
| `COPY go.mod go.sum` | Копирование файлов зависимостей |
| `RUN go mod download` | Скачивание зависимостей (кешируется) |
| `COPY . .` | Копирование исходного кода |
| `RUN go build` | Компиляция статического бинарника |

**Особенности:**
- `CGO_ENABLED=0` - отключает CGO для статической сборки
- `GOOS=linux` - явное указание целевой ОС
- Бинарник собирается в `/app/bin/`

**Stage 2: Runner**
| Шаг | Назначение |
|-----|------------|
| `FROM alpine:latest` | Минимальный базовый образ (~5MB) |
| `RUN addgroup/adduser` | Создание непривилегированного пользователя |
| `COPY --from=builder` | Копирование только бинарника из builder стадии |
| `USER` | Запуск от непривилегированного пользователя |
| `EXPOSE` | Документирование портов |
| `CMD` | Команда запуска |

**Преимущества multi-stage:**
- **Маленький размер** - итоговый образ содержит только бинарник (~15MB вместо ~800MB)
- **Безопасность** - нет компиляторов и лишних утилит
- **Кеширование** - зависимости скачиваются только при изменении go.mod
- **Воспроизводимость** - одинаковые бинарники при каждой сборке

### Практическое занятие №8 (CI/CD Pipeline)

#### Выбор платформы: GitHub Actions
- Репозиторий уже размещен на GitHub
- Тесная интеграция с GitHub Container Registry (GHCR)
- Бесплатный для публичных репозиториев
- Простая настройка через YAML файлы
- Встроенные возможности для Docker сборки

#### Структура pipeline

**Job 1: Test (тестирование)**
| Шаг | Назначение |
|-----|------------|
| Checkout | Клонирование репозитория |
| Setup Go | Установка Go 1.25 с кешированием |
| Install dependencies | `go mod download` |
| go vet | Статический анализ кода |
| go test | Запуск тестов с флагами `-race -cover` |
| Upload coverage | Отправка отчета о покрытии в Codecov |

**Job 2: Build (сборка бинарников)**
| Шаг | Назначение |
|-----|------------|
| Checkout | Клонирование репозитория |
| Setup Go | Установка Go с кешированием |
| Build Auth | Компиляция Auth сервиса |
| Build Tasks | Компиляция Tasks сервиса |
| Upload artifacts | Сохранение бинарников |

**Job 3: Docker (сборка и публикация образов)**
| Шаг | Назначение |
|-----|------------|
| Checkout | Клонирование репозитория |
| Setup Buildx | Настройка Docker Buildx |
| Login to GHCR | Аутентификация в GitHub Container Registry |
| Extract metadata | Генерация тегов для образов |
| Build and push Auth | Сборка и публикация Auth образа |
| Build and push Tasks | Сборка и публикация Tasks образа |

**Job 4: Deploy (деплой на сервер)**
| Шаг | Назначение |
|-----|------------|
| Checkout | Клонирование репозитория |
| Setup SSH | Настройка SSH агента с приватным ключом |
| Add host to known_hosts | Добавление сервера в known_hosts |
| Create .env | Создание файла с переменными окружения |
| Deploy with Docker Compose | Копирование файлов и запуск контейнеров |

**Job 5: Notify (уведомление)**
| Шаг | Назначение |
|-----|------------|
| Send Telegram notification | Отправка результатов pipeline в Telegram |

#### Версионирование Docker-образов
```yaml
tags: |
  type=sha,format=short    # ghcr.io/username/repo/auth:sha-a1b2c3d
  type=ref,event=branch     # ghcr.io/username/repo/auth:main
  type=raw,value=latest,enable={{is_default_branch}}  # ghcr.io/username/repo/auth:latest
```
#### Примеры тегов:
- ghcr.io/username/auth:latest
- ghcr.io/username/auth:main
- ghcr.io/username/auth:sha-a1b2c3d
- ghcr.io/username/tasks:latest
- ghcr.io/username/tasks:sha-f4e5d6c

### Практическое занятие №9 (Распределённый кэш Redis)
#### Реализованные возможности:
- **Cache-aside pattern** для GET /v1/tasks/{id} и GET /v1/tasks
- **TTL с jitter** (120s + 0-30s) для предотвращения cache avalanche
- **Инвалидация кэша** при создании/обновлении/удалении задач
- **Деградация** - при недоступности Redis сервис продолжает работать через БД
#### Ключи кэша:
| Тип | Ключ | TTL | Инвалидация |
|-----|------|-----|-------------|
| Задача | `tasks:task:{id}` | 120-150s | При update/delete |
| Список | `tasks:list:{subject}` | 120-150s | При create/update/delete |

#### Переменные окружения для Redis:
| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `REDIS_ADDR` | redis:6379 | Адрес Redis сервера |
| `REDIS_PASSWORD` | - | Пароль Redis |
| `REDIS_DB` | 0 | Номер базы данных |
| `CACHE_TTL_SECONDS` | 120 | Базовый TTL в секундах |
| `CACHE_TTL_JITTER_SECONDS` | 30 | Максимальный jitter |

### Практическое занятие №10 (Горизонтальное масштабирование: использование Load Balancer (NGINX))
#### Реализованные возможности:
- **Горизонтальное масштабирование** сервиса tasks с запуском 3 реплик
- **Load Balancing** через NGINX с распределением трафика между репликами
- **Идентификация инстанса** через заголовок X-Instance-ID для отладки
- **Единая точка входа** — доступ ко всем репликам через порт 8080
- **Отказоустойчивость** — при падении одной реплики трафик перенаправляется на остальные

#### Конфигурация реплик:
| Параметр | tasks_1 | tasks_2 | tasks_3 |
| :--- | :--- | :--- | :--- |
| INSTANCE_ID | tasks-1 | tasks-2 | tasks-3 |
| Внутренний порт | 8082 | 8082 | 8082 |
| Внешний порт | — | — | — |
| Доступ через LB | 8080 | 8080 | 8080 |

#### Переменные окружения для масштабирования:
| Переменная | По умолчанию | Описание |
| :--- | :--- | :--- |
| `INSTANCE_ID` | tasks-1 | Уникальный идентификатор инстанса |
| `INTERNAL_PORT` | 8082 | Внутренний порт сервиса |
| `LB_PORT` | 8080 | Порт NGINX Load Balancer |
| `NGINX_WORKER_CONNECTIONS` | 1024 | Максимальное количество соединений NGINX |
| `UPSTREAM_HEALTH_CHECK` | enabled | Проверка здоровья upstream серверов |

#### Архитектура развёртывания:
| Компонент | Количество | Назначение |
| :--- | :--- | :--- |
| tasks реплики | 3 | Обработка запросов |
| NGINX | 1 | Распределение трафика |
| Redis | 1 | Кэширование (занятие №9) |
| PostgreSQL | 1 | Хранение данных |

## Практическое занятие №11 (GraphQL API)
### GraphQL Service
- **Порт**: 8090
- **Endpoint**: `/query`
- **Playground**: `http://193.233.175.221:8090/`
- **Технологии**: gqlgen, GraphQL, PostgreSQL

### Схема GraphQL

```graphql
type Task {
  id: ID!
  title: String!
  description: String
  due_date: String
  done: Boolean!
  created_at: String
  updated_at: String
}

input CreateTaskInput {
  title: String!
  description: String
  due_date: String
}

input UpdateTaskInput {
  title: String
  description: String
  due_date: String
  done: Boolean
}

type Query {
  tasks: [Task!]!
  task(id: ID!): Task
}

type Mutation {
  createTask(input: CreateTaskInput!): Task!
  updateTask(id: ID!, input: UpdateTaskInput!): Task!
  deleteTask(id: ID!): Boolean!
}
```
#### Примеры запросов
##### Получить все задачи
```graphql
query {
  tasks {
    id
    title
    done
  }
}
```
##### Получить задачу по ID
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
##### Создать задачу
```graphql
mutation CreateTask($input: CreateTaskInput!) {
  createTask(input: $input) {
    id
    title
    done
  }
}
```
##### Обновить задачу
```graphql
mutation UpdateTask($id: ID!, $input: UpdateTaskInput!) {
  updateTask(id: $id, input: $input) {
    id
    title
    done
  }
}
```
##### Удалить задачу
```graphql
mutation DeleteTask($id: ID!) {
  deleteTask(id: $id)
}
```
## Практическое занятие №12 (Сравнение REST и GraphQL)
### Выбранный UI-сценарий
#### Экран списка задач
- Нужны поля: id, title, done
#### Экран деталей задачи
- Нужны поля: id, title, description, done
#### Действия
- Создать задачу
- Отметить задачу выполненной (обновить)

### Реализация через REST
#### Получить список задач
Размер ответа: ~180 байт (с 2 задачами)
#### Получить детали задачи
Размер ответа: ~220 байт

### Реализация через GraphQL
#### Получить список задач (только нужные поля)
Размер ответа: ~160 байт (с 2 задачами) - на 12% меньше REST
#### Получить детали задачи (только нужные поля)
Размер ответа: ~120 байт - на 45% меньше REST

#### Основные отличия
| Критерий | REST | GraphQL |
| :--- | :--- | :--- |
| Количество запросов для сценария | 2 (список + детали) | 1 (можно получить всё в одном запросе) |
| Объём данных (список, 2 задачи) | 180 байт | 160 байт |
| Объём данных (детали) | 220 байт | 120 байт |
| Over-fetching | Есть (получаем лишние поля) | Нет (только запрошенные поля) |
| Under-fetching | Нет | Нет |
| Обработка ошибок | HTTP статусы (200, 400, 401, 404, 500) | HTTP 200 + поле errors |
| Кэширование | Простое (по URL) | Сложное (один endpoint) |
| Документация | Swagger/OpenAPI (встроенная) | GraphQL Schema (самодокументируемая) |
| Версионирование | Через URL (/v1/, /v2/) | Эволюция схемы |

### Итоговый вывод
#### Когда выбирать REST
- Простые API - CRUD операции над ресурсами
- Кэширование критично - CDN, браузерное кэширование
- Стандартизация - чёткие HTTP статусы и методы
- Микросервисы - простота интеграции
- Публичные API - широкая поддержка инструментов
#### Когда выбирать GraphQL:
- Сложные клиенты - мобильные приложения с ограниченным трафиком
- Агрегация данных - нужно объединять данные из нескольких источников
- Быстрая эволюция - клиенты сами выбирают нужные поля
- Over-fetching проблема - нужно минимизировать передаваемые данные
- Разные клиенты - веб, мобильные, десктоп с разными потребностями

## Практическое занятие №13 (Подключение к RabbitMQ. Отправка и получение сообщений)
### Режим публикации: "best effort"
- Если RabbitMQ недоступен - логируем ошибку, но задача создаётся
- Асинхронная публикация (не блокирует ответ клиенту)
- Persistent сообщения (переживают рестарт RabbitMQ)

### Устройство worker
#### Consumer
- Подключается к RabbitMQ
- Объявляет ту же очередь (durable)
- Prefetch = 1 (обрабатывает по одному сообщению)
- Manual ack (подтверждение после обработки)
- Логирует полученные события
#### 2 экземпляра worker
- worker-1 и worker-2
- Распределяют нагрузку (round-robin)
- Каждый подтверждает свои сообщения

### Практическое занятие №14 (Реализация очереди задач (producer–consumer): retries, DLQ, идемпотентность)

#### Реализованные возможности:
- **Асинхронная обработка** задач через очередь сообщений
- **Механизм Retry** с ограниченным количеством попыток и задержкой
- **Dead Letter Queue (DLQ)** для хранения неудачных сообщений
- **Идемпотентность** обработки сообщений для предотвращения дубликатов

#### Конфигурация очередей:
| Очередь | Назначение | Durable | Особенности |
| :--- | :--- | :--- | :--- |
| `task_jobs` | Основная очередь задач | Да | Получает задания от producer |
| `task_jobs_retry` | Очередь повторных попыток | Да | TTL = 10 секунд, возврат в основную очередь |
| `task_jobs_dlq` | Dead Letter Queue | Да | Хранит сообщения, превысившие лимит попыток |

#### Настройка Dead Letter Exchange (DLX):
| Очередь | x-dead-letter-exchange | Routing Key | Поведение |
| :--- | :--- | :--- | :--- |
| `task_jobs` | `task_jobs_dlx` | `dlq` | При превышении попыток → DLQ |
| `task_jobs_retry` | `""` (default) | — | Возврат в основную очередь через TTL |

#### Формат сообщения (Job):
| Поле | Тип | Описание |
| :--- | :--- | :--- |
| `job_type` | string | Тип задания (process_task) |
| `task_id` | string | Идентификатор задачи для обработки |
| `attempt` | int | Номер попытки (1, 2, 3...) |
| `message_id` | string | Уникальный UUID для идемпотентности |
| `created_at` | string | Временная метка создания |

#### Политика повторных попыток (Retry Policy):
| Параметр | Значение |
| :--- | :--- |
| Максимальное количество попыток | 3 |
| Задержка между попытками | 10 секунд (через retry queue) |
| Что считается ошибкой | Любая ошибка при обработке + ошибки для task_id содержащих "fail" |
| Действие при превышении попыток | Отправка в DLQ |

#### Механизм Retry:
1. Worker получает сообщение и начинает обработку
2. При ошибке увеличивается поле `attempt` на 1
3. Если `attempt <= 3` → публикация в `task_jobs_retry` с TTL 10 сек
4. Через 10 секунд сообщение возвращается в основную очередь `task_jobs`
5. Если `attempt > 3` → отправка в Dead Letter Queue

#### Идемпотентность:
| Параметр | Значение |
| :--- | :--- |
| Ключ идемпотентности | `message_id` (UUID) |
| Где хранятся обработанные ID | В памяти (`storage.ProcessedMessages`) |
| TTL хранения | 24 часа |
| Проверка | Перед обработкой проверяется наличие message_id в хранилище |

### Практическое занятие №15 (Деплой приложения на VPS. Настройка system)
#### Подготовка VPS (на сервере)
```
apt update && apt upgrade -y

# Создать пользователя для сервиса
useradd --system --no-create-home --shell /usr/sbin/nologin tasksuser

# Проверка
id tasksuser
# uid=998(tasksuser) gid=996(tasksuser) groups=996(tasksuser)

# Директория для бинарника
mkdir -p /opt/tasks/bin

# Директория для конфигов
mkdir -p /etc/tasks

# Директория для логов
mkdir -p /var/log/tasks

# Назначение владельца
chown -R tasksuser:tasksuser /opt/tasks
chown -R tasksuser:tasksuser /var/log/tasks
chmod 755 /opt/tasks
chmod 755 /opt/tasks/bin
```
#### Сборка бинарника (на локальной машине)
```
# Переходим в корень проекта
cd /путь/к/tech-ip-sem2

# Собираем бинарник для Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o tasks-bin ./services/tasks/cmd/tasks; $env:GOOS="windows"; $env:GOARCH="amd64"

# Проверяем, что создан
Get-ChildItem tasks-bin
Ответ:
    Каталог: C:\Users\User\Downloads\tech-ip-sem2


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----        09.03.2026     19:43       22210702 tasks-bin
```
#### Копирование на сервер (локально)
```
# Копирование бинарника
scp tasks-bin root@193.233.175.221:/tmp/tasks-bin

# Копируем необходимых файлов конфигурации
scp .env root@193.233.175.221:/tmp/tasks.env
```
#### Размещение на сервере (на сервере)
```
# Перемещение бинарника
mv /tmp/tasks-bin /opt/tasks/bin/tasks
chmod 755 /opt/tasks/bin/tasks
chown tasksuser:tasksuser /opt/tasks/bin/tasks

# Создание конфига с переменными окружения
cat > /etc/tasks/tasks.env << 'EOF'
# Tasks Service Configuration
TASKS_PORT=9082  # Другой порт, чтобы не конфликтовать с Docker
AUTH_GRPC_ADDR=localhost:50051
DB_HOST=localhost
DB_PORT=5432
DB_USER=appuser
DB_PASSWORD=AppP@ssw0rd!2026
DB_NAME=tasksdb
DB_SSLMODE=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL_SECONDS=120
CACHE_TTL_JITTER_SECONDS=30
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE=task_events
INSTANCE_ID=tasks-systemd
LOG_LEVEL=info
EOF

# Выставляем правильные права
chown root:root /etc/tasks/tasks.env
chmod 600 /etc/tasks/tasks.env
```
#### Systemd unit файл
```
cat > /etc/systemd/system/tasks.service << 'EOF'
[Unit]
Description=Tasks Service (systemd demo)
After=network.target
Wants=network.target

[Service]
Type=simple
User=tasksuser
Group=tasksuser
WorkingDirectory=/opt/tasks

# Переменные окружения из файла
EnvironmentFile=/etc/tasks/tasks.env

# Запуск бинарника
ExecStart=/opt/tasks/bin/tasks
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5

# Безопасность
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
ReadWritePaths=/var/log/tasks

# Лимиты
LimitNOFILE=65535
LimitNPROC=65535

# Логирование
StandardOutput=journal
StandardError=journal
SyslogIdentifier=tasks

[Install]
WantedBy=multi-user.target
EOF
```
#### Параметры
Параметр	Значение	Пояснение
Type=simple	-	Процесс запускается и остается в памяти
User=tasksuser	-	Запуск от непривилегированного пользователя
EnvironmentFile	/etc/tasks/tasks.env	Конфиг отдельно от кода
Restart=always	-	Перезапуск при любом завершении (кроме stop)
RestartSec=5	5 секунд	Пауза перед перезапуском
NoNewPrivileges=true	-	Запрет на повышение привилегий
ProtectSystem=strict	-	Только чтение системных директорий
ProtectHome=true	-	Нет доступа к /home
PrivateTmp=true	-	Свой изолированный /tmp
#### Запуск и управление сервисом
```
# Перечитываем конфиги systemd
systemctl daemon-reload

# Запускаем сервис
systemctl start tasks

# Включаем автозапуск
systemctl enable tasks

# Проверяем статус
systemctl status tasks
```
#### Основные команды systemd
```
# Статус
systemctl status tasks

# Стоп
systemctl stop tasks

# Старт
systemctl start tasks

# Рестарт
systemctl restart tasks

# Перезагрузка конфига (без остановки)
systemctl reload tasks

# Проверка включен ли автозапуск
systemctl is-enabled tasks

# Отключение автозапуска
systemctl disable tasks
```
#### Просмотр логов через journalctl
```
# Последние 50 логов
journalctl -u tasks -n 50 --no-pager

# Логи в реальном времени
journalctl -u tasks -f

# Логи за последний час
journalctl -u tasks --since "1 hour ago"

# Логи с определенного времени
journalctl -u tasks --since "2026-03-10 10:00:00"

# Логи с определенного времени до сейчас
journalctl -u tasks --since "2026-03-10 10:00:00" -f

# Логи в JSON формате
journalctl -u tasks -o json-pretty
```
### Практическое занятие №16 (Публикация приложения в Kubernetes)
#### Kubernetes стенд
- Использован: Minikube v1.38.1 на Windows 11 Pro

#### Ключевые параметры Deployment
- replicas: 2 - запущено 2 экземпляра сервиса
- image: ghcr.io/mamuer/technology_2_sem/tasks:latest - используемый образ
- containerPort: 8082 - порт контейнера
- envFrom - подключение ConfigMap и Secret
- livenessProbe - проверка живости (перезапуск при падении)
- readinessProbe - проверка готовности (трафик только когда готов)

## Команды запуска и сборки
- make graphql-run        # Запустить локально
- make graphql-build      # Собрать Docker образ
- make graphql-up         # Запустить в Docker
- make graphql-logs       # Посмотреть логи
- make graphql-test-query # Протестировать запрос
- make graphql-test-create # Протестировать создание
### Основные команды
- make check Проверка кода и форматирования
- make tree Показать структуру проекта
- make help Показать справку
- make generate Сгенерировать код из proto файлов
#### Docker команды
- make docker-build	Собрать образы
- make docker-up	Запустить контейнеры
- make docker-down	Остановить контейнеры
- make docker-reset	Полный сброс (clean + rebuild)
- make docker-logs	Просмотр логов всех сервисов
- make docker-logs-auth	Логи только Auth сервиса
- make docker-logs-tasks	Логи только Tasks сервиса
- make docker-ps	Статус контейнеров
- make docker-reset Полный сброс (clean + rebuild)
### Мониторинг
- make monitor-up	Запустить Prometheus + Grafana
- make monitor-down	Остановить мониторинг
- make monitor-logs	Логи мониторинга
### HTTPS и безопасность
- make gen-cert	Сгенерировать SSL сертификаты
- make p22-up	Запустить HTTPS + PostgreSQL стек
- make p22-down	Остановить HTTPS стек
### Тестирование
- make test-load	Сгенерировать тестовую нагрузку
- make test-load-create	Создать тестовые задачи
- make fast-auth	Запустить Auth локально
- make fast-tasks	Запустить Tasks локально
## Сборка образов
### Сборка Auth сервиса
```
cd services/auth
docker build -t techip-auth:1.0 .
```
### Сборка Tasks сервиса
```
cd services/tasks
docker build -t techip-tasks:1.0 .
```
### Сборка из корня проекта
- docker build -f services/auth/Dockerfile -t techip-auth:1.0 .
- docker build -f services/tasks/Dockerfile -t techip-tasks:1.0 .
## Запуск отдельных контейнеров
## Запуск Auth (с пробросом портов)
```
docker run -d \
  --name auth-service \
  -p 8081:8081 \
  -p 50051:50051 \
  -e AUTH_PORT=8081 \
  -e AUTH_GRPC_PORT=50051 \
  techip-auth:1.0
```
## Запуск Tasks (с подключением к Auth)
```
docker run -d \
  --name tasks-service \
  -p 8082:8082 \
  -e TASKS_PORT=8082 \
  -e AUTH_GRPC_ADDR=host.docker.internal:50051 \
  techip-tasks:1.0
```
## Полезные команды Docker
- docker images	Просмотр образов
- docker ps	Просмотр запущенных контейнеров
- docker ps -a	Просмотр всех контейнеров
- docker logs auth-service	Логи контейнера auth
- docker logs tasks-service	Логи контейнера tasks
- docker stop auth-service	Остановка контейнера
- docker rm auth-service	Удаление контейнера
- docker rmi techip-auth:1.0	Удаление образа
- docker exec -it auth-service sh	Интерактивный вход в контейнер
## Передача переменных при запуске
### Через командную строку
- docker run -e AUTH_PORT=9090 -e AUTH_GRPC_PORT=50051 techip-auth:1.0
### Через файл .env
- docker run --env-file .env techip-tasks:1.0
## Команды Docker Compose
- docker-compose up -d	Запуск всех сервисов
- docker-compose up -d --build	Запуск с пересборкой образов
- docker-compose down	Остановка всех сервисов
- docker-compose down -v	Остановка с удалением томов
- docker-compose logs -f	Просмотр логов
- docker-compose ps	Просмотр статуса
- docker-compose restart tasks	Перезапуск конкретного сервиса
- docker-compose up -d --scale tasks=3	Масштабирование сервиса

## Описание
### gRPC коммуникация между сервисами
- **Auth service**: работает как gRPC сервер на порту 50051 (параллельно с HTTP)
- **Tasks service**: использует gRPC клиент для проверки токенов
- **Deadline**: каждый вызов имеет таймаут 3 секунды
- **Proto контракт**: описан в `proto/auth.proto`

### Выбранный логгер: **Uber Zap**
- Структурированные логи в JSON формате по умолчанию
- Поддержка уровней логирования (DEBUG, INFO, WARN, ERROR)
- Возможность добавления контекстных полей

### TLS терминирование: **NGINX**
- Соответствие индустриальным стандартам
- Разделение ответственности
- Централизованное управление сертификатами
- Дополнительные возможности: балансировка, кеширование, сжатие
- Безопасность: проверенная кодовая база для работы с TLS

### Стандарт полей логов
**Обязательные поля:**
- `level` - уровень логирования (debug/info/warn/error)
- `ts` - временная метка в ISO8601 формате
- `service` - имя сервиса (auth/tasks)
- `request_id` - идентификатор запроса для трассировки
- `method` - HTTP метод
- `path` - путь запроса
- `status` - HTTP код ответа
- `duration_ms` - время обработки в миллисекундах

**Для ошибок дополнительно:**
- `error` - текст ошибки (без чувствительных данных)
- `component` - компонент, где произошла ошибка (repository, handler, auth_client)

### Реализация request-id
- Request-id извлекается из заголовка `X-Request-ID` или генерируется новый (UUID)
- Добавляется в контекст запроса и в ответ (заголовок `X-Request-ID`)
- Прокидывается в gRPC вызовы через метаданные
- Позволяет трассировать запрос через оба сервиса

### Защита от CSRF (Double Submit Cookie)
**Как работает:**
1. При логине сервер устанавливает две cookies:
   - `session_id` (HttpOnly, Secure, SameSite=Lax) - для аутентификации
   - `csrf_token` (не HttpOnly, Secure, SameSite=Lax) - для CSRF защиты
2. Клиент читает `csrf_token` из cookie и отправляет его в заголовке `X-CSRF-Token`
3. Сервер сравнивает токен из cookie и из заголовка
4. При несовпадении - возвращает 403 Forbidden

**Используемые cookies:**

| Cookie | HttpOnly | Secure | SameSite | Max-Age | Назначение |
|--------|----------|--------|----------|---------|------------|
| session_id | Да | Да | Lax | 86400 (24ч) | Идентификатор сессии |
| csrf_token | Нет | Да | Lax | 86400 (24ч) | Токен для CSRF защиты |

**Защищаемые методы:**
- POST
- PUT
- PATCH
- DELETE

### Защита от XSS
**Реализованные меры:**
1. **Санитизация ввода:**
   - Функция `SanitizeText()` - экранирует HTML специальные символы
   - Функция `SanitizeHTML()` - удаляет HTML теги и экранирует символы
   - Функция `ValidateAndSanitizeDescription()` - проверка длины и очистка

2. **Заголовки безопасности:**
   - `X-Content-Type-Options: nosniff` - защита от MIME sniffing
   - `X-Frame-Options: DENY` - защита от clickjacking
   - `X-XSS-Protection: 1; mode=block` - включение XSS фильтра в браузерах
   - `Content-Security-Policy: default-src 'self'` - ограничение источников контента
   - `Referrer-Policy: strict-origin-when-cross-origin` - контроль Referrer
   - `Strict-Transport-Security: max-age=31536000; includeSubDomains; preload` - HSTS

### Примеры логов
**Успешный запрос (auth):**
```json
{
  "level": "info",
  "ts": "2026-02-24T10:30:45.123Z",
  "service": "auth",
  "request_id": "test-request-001",
  "method": "POST",
  "path": "/v1/auth/login",
  "status": 200,
  "duration_ms": 45,
  "remote_ip": "172.18.0.1",
  "user_agent": "curl/7.68.0"
}
```
**Запрос с ошибкой (tasks):**

```json
{
  "level": "warn",
  "ts": "2026-02-24T10:31:22.456Z",
  "service": "tasks",
  "request_id": "test-request-002",
  "method": "GET",
  "path": "/v1/tasks",
  "status": 401,
  "duration_ms": 12,
  "remote_ip": "172.18.0.1",
  "user_agent": "curl/7.68.0",
  "error": "unauthorized"
}
```
**Межсервисный вызов (auth gRPC):**

```json
{
  "level": "info",
  "ts": "2026-02-24T10:32:05.789Z",
  "service": "auth",
  "request_id": "test-request-003",
  "method": "gRPC",
  "path": "/auth.AuthService/Verify",
  "duration_ms": 23,
  "subject": "student"
}
```
### CI/CD Pipeline 
#### Файл конфигурации
.github/workflows/ci.yml

## Структура проекта
```
C:.
│   .dockerignore
│   .env
│   .env.example
│   .gitignore
│   docker-compose.override.yml
│   docker-compose.prod.yml
│   go.mod
│   go.sum
│   Makefile
│   passwords.txt
│   README.md
│   tasks-bin
│
├───.github
│   └───workflows
│           ci.yml
│
├───deploy
│   ├───lb
│   │       nginx.conf
│   │
│   ├───monitoring
│   │   │   docker-compose.yml
│   │   │   prometheus.yml
│   │   │
│   │   └───grafana
│   │       └───provisioning
│   │           ├───dashboards
│   │           │       dashboard.yml
│   │           │       tasks-dashboard.json
│   │           │
│   │           └───datasources
│   │                   prometheus.yml
│   │
│   ├───redis
│   │       docker-compose.yml
│   │
│   └───tls
│       │   docker-compose.yml
│       │   generate-cert.sh
│       │   init.sql
│       │   nginx.conf
│       │
│       └───certs
│               cert.pem
│               key.pem
│
├───docs
│       pz17_diagram.md
│       pz20_metrics.md
│       pz_api.md
│
├───PR1-16
│
├───proto
│   │   auth.proto
│   │
│   └───gen
│       └───go
│           └───auth
│                   auth.pb.go
│                   auth_grpc.pb.go
│
├───services
│   ├───auth
│   │   │   Dockerfile
│   │   │
│   │   ├───cmd
│   │   │   └───auth
│   │   │           main.go
│   │   │
│   │   └───internal
│   │       ├───config
│   │       │       config.go
│   │       │
│   │       ├───grpc
│   │       │       server.go
│   │       │
│   │       ├───http
│   │       │       handler.go
│   │       │
│   │       ├───models
│   │       │       session.go
│   │       │
│   │       └───service
│   │               auth.go
│   │               auth_test.go
│   │               session_service.go
│   │
│   ├───graphql
│   │   │   Dockerfile
│   │   │   gqlgen.yml
│   │   │
│   │   ├───cmd
│   │   │   └───graphql
│   │   │           main.go
│   │   │
│   │   ├───graph
│   │   │   │   schema.graphqls
│   │   │   │
│   │   │   ├───generated
│   │   │   │       generated.go
│   │   │   │
│   │   │   ├───model
│   │   │   │       models_gen.go
│   │   │   │
│   │   │   └───resolvers
│   │   │           resolver.go
│   │   │           schema.resolvers.go
│   │   │
│   │   ├───internal
│   │   │   ├───middleware
│   │   │   │       auth.go
│   │   │   │
│   │   │   ├───repository
│   │   │   │       task_repository.go
│   │   │   │
│   │   │   └───service
│   │   │           task_service.go
│   │   │
│   │   └───tools
│   │           tools.go
│   │
│   ├───tasks
│   │   │   Dockerfile
│   │   │
│   │   ├───cmd
│   │   │   └───tasks
│   │   │           main.go
│   │   │
│   │   └───internal
│   │       ├───cache
│   │       │       redis_cache.go
│   │       │
│   │       ├───client
│   │       │   └───authclient
│   │       │           client.go
│   │       │
│   │       ├───config
│   │       │       config.go
│   │       │
│   │       ├───http
│   │       │       handlers.go
│   │       │       job_handlers.go
│   │       │
│   │       ├───jobs
│   │       │       models.go
│   │       │
│   │       ├───models
│   │       │       task.go
│   │       │
│   │       ├───rabbitmq
│   │       │       job_publisher.go
│   │       │       publisher.go
│   │       │
│   │       ├───repository
│   │       │       task_repository.go
│   │       │
│   │       └───service
│   │               tasks.go
│   │               tasks_test.go
│   │
│   └───worker
│       │   Dockerfile
│       │
│       ├───cmd
│       │   └───worker
│       │           main.go
│       │
│       └───internal
│           ├───consumer
│           │       consumer.go
│           │       job_consumer.go
│           │
│           ├───models
│           │       event.go
│           │       job.go
│           │
│           ├───processor
│           │       task_processor.go
│           │
│           └───storage
│                   memory.go
│
└───shared
    ├───cookies
    │       cookie_helper.go
    │
    ├───httpx
    │       client.go
    │
    ├───logger
    │       logger.go
    │
    ├───metrics
    │       metrics.go
    │
    ├───middleware
    │       accesslog.go
    │       csrf.go
    │       debug.go
    │       logging.go
    │       metrics.go
    │       requestid.go
    │       security_headers.go
    │
    └───sanitize
            sanitizer.go
```
## Скриншоты работы проекта
### Запуск докер контейнеров
![фото1](./PR1-16/Screenshot_1.png)
### Получение токена доступа
#### Удачное
![фото2](./PR1-16/Screenshot_2.png)
#### Ошибка авторизации
![фото3](./PR1-16/Screenshot_3.png)
#### Ошибка формата
![фото4](./PR1-16/Screenshot_4.png)
### Проверка валидности токена
#### Удачное
![фото5](./PR1-16/Screenshot_5.png)
#### Ошибка авторизации
![фото6](./PR1-16/Screenshot_6.png)
### Создание новой задачи
![фото7](./PR1-16/Screenshot_7.png)
### Получение списка всех задач
![фото8](./PR1-16/Screenshot_8.png)
### Получение задачи по ID
#### Удачное
![фото9](./PR1-16/Screenshot_9.png)
#### Не найдено задание
![фото10](./PR1-16/Screenshot_10.png)
### Обновление задачи
![фото11](./PR1-16/Screenshot_11.png)
### Удаление задачи
![фото12](./PR1-16/Screenshot_12.png)
### Проверка отказоустойчивости системы
#### Verify пошёл через gRPC
![фото13](./PR1-16/Screenshot_13.png)
#### Остановка Auth
![фото14](./PR1-16/Screenshot_14.png)
#### Ошибка после остановки Auth
![фото15](./PR1-16/Screenshot_15.png)
#### Ошибка таймаут Auth
![фото16](./PR1-16/Screenshot_16.png)
### Логи после авторизации
![фото17](./PR1-16/Screenshot_17.png)
### Логи после ошибки при авторизации
![фото18](./PR1-16/Screenshot_18.png)
### Логи после запроса списка задач
![фото19](./PR1-16/Screenshot_19.png)
### Метрики
#### Tasks сервиса
![фото20](./PR1-16/Screenshot_20.png)
#### Auth сервиса
![фото21](./PR1-16/Screenshot_21.png)
### Prometheus
#### RPS (Requests Per Second)
![фото22](./PR1-16/Screenshot_22.png)
#### Количество ошибок 4xx
![фото23](./PR1-16/Screenshot_23.png)
#### 95-й перцентиль длительности запросов
![фото24](./PR1-16/Screenshot_24.png)
#### Текущие активные запросы
![фото25](./PR1-16/Screenshot_25.png)
#### Статус сервисов
![фото26](./PR1-16/Screenshot_26.png)
### Grafana
#### Data source connection
![фото27](./PR1-16/Screenshot_27.png)
#### График RPS (Requests Per Second)
![фото28](./PR1-16/Screenshot_28.png)
#### График ошибок (Errors)
![фото29](./PR1-16/Screenshot_29.png)
#### График задержек (p95 Latency)
![фото30](./PR1-16/Screenshot_30.png)
#### График активных запросов (In-flight)
![фото31](./PR1-16/Screenshot_31.png)
### Генерация сертификатов
![фото32](./PR1-16/Screenshot_32.png)
### Ошибка проверки состояния сервиса
![фото33](./PR1-16/Screenshot_33.png)
### HTTPS
#### Состояние
![фото34](./PR1-16/Screenshot_34.png)
#### Метрики
![фото35](./PR1-16/Screenshot_35.png)
#### Создание задачи
![фото36](./PR1-16/Screenshot_36.png)
#### Поиск
![фото37](./PR1-16/Screenshot_37.png)
#### Поиск с SQL-инъекцией
![фото38](./PR1-16/Screenshot_38.png)
#### Поиск с SQL-инъекцией (не работает)
![фото39](./PR1-16/Screenshot_39.png)
### CSRF
#### Login
![фото40](./PR1-16/Screenshot_40.png)
#### Cookies
![фото41](./PR1-16/Screenshot_41.png)
#### Отсутствие CSRF
![фото42](./PR1-16/Screenshot_42.png)
#### Создание задачи с CSRF
![фото43](./PR1-16/Screenshot_43.png)
#### Проверка работы защиты
![фото44](./PR1-16/Screenshot_44.png)
#### Проверка сетевого взаимодействия
![фото45](./PR1-16/Screenshot_45.png)

![фото46](./PR1-16/Screenshot_46.png)

![фото47](./PR1-16/Screenshot_47.png)

![фото48](./PR1-16/Screenshot_48.png)
### Успешный прогон Actions
![фото49](./PR1-16/Screenshot_49.png)
### Переменные Actions на GitHub
![фото50](./PR1-16/Screenshot_50.png)
### Телеграм бот с уведомлением о состоянии jobs
![фото51](./PR1-16/Screenshot_51.png)
### Ключи кэша (Redis)
![фото52](./PR1-16/Screenshot_52.png)
### Логи с hit/miss
![фото53](./PR1-16/Screenshot_53.png)
### Сравнение времени запросов (hit vs miss)
![фото54](./PR1-16/Screenshot_54.png)
### TTL ключа
![фото55](./PR1-16/Screenshot_55.png)
### Проверка деградации (остановка Redis)
![фото56](./PR1-16/Screenshot_56.png)
### Конфигурация Redis в docker-compose
![фото57](./PR1-16/Screenshot_57.png)
### Переменные окружения
![фото58](./PR1-16/Screenshot_58.png)
### Одновременная работа 3х реплик
![фото59](./PR1-16/Screenshot_59.png)
### Запросы обрабатываются разными репликами
![фото60](./PR1-16/Screenshot_60.png)
### Уведомление о контейнирах
![фото61](./PR1-16/Screenshot_61.png)
### Доступ к GraphQL Playground
![фото62](./PR1-16/Screenshot_62.png)
### GraphQL Playground с запросом tasks
![фото63](./PR1-16/Screenshot_63.png)
### GraphQL Playground с мутацией createTask
![фото64](./PR1-16/Screenshot_64.png)
### GraphQL Playground с запросом task(id)
![фото65](./PR1-16/Screenshot_65.png)
### GraphQL Playground с мутацией updateTask
![фото66](./PR1-16/Screenshot_66.png)
### GraphQL Playground с мутацией deleteTask
![фото67](./PR1-16/Screenshot_67.png)
### RabbitMQ Queues и очередь task_events
![фото68](./PR1-16/Screenshot_68.png)
### Создание задачи и логи с подтверждением публикации
![фото69](./PR1-16/Screenshot_69.png)
### Создание normal-task
![фото70](./PR1-16/Screenshot_70.png)
### Лог normal-task
![фото72](./PR1-16/Screenshot_72.png)
### Создание fail-task
![фото71](./PR1-16/Screenshot_71.png)
### Лог fail-task
![фото73](./PR1-16/Screenshot_73.png)
### DLQ
![фото74](./PR1-16/Screenshot_74.png)
### RabbitMQ UI с очередями
![фото75](./PR1-16/Screenshot_75.png)
### Сборка бинарника
![фото76](./PR1-16/Screenshot_76.png)
### Создание пользователя и директорий
![фото77](./PR1-16/Screenshot_77.png)
### Копирование и размещение на сервере
![фото78](./PR1-16/Screenshot_78.png)
### Systemd unit файл
![фото79](./PR1-16/Screenshot_79.png)
### Запуск и управление сервисом
![фото80](./PR1-16/Screenshot_80.png)
### Проверка health endpoint
![фото81](./PR1-16/Screenshot_81.png)
### Вывод активации сервиса
![фото82](./PR1-16/Screenshot_82.png)
### Работа метода tasks
![фото83](./PR1-16/Screenshot_83.png)
### k8s-services
![фото84](./PR1-16/Screenshot_84.png)
### k8s-deployment
![фото85](./PR1-16/Screenshot_85.png)
### k8s-logs
![фото86](./PR1-16/Screenshot_86.png)
### k8s-pods
![фото87](./PR1-16/Screenshot_87.png)
### port-forward
![фото88](./PR1-16/Screenshot_88.png)
### health
![фото89](./PR1-16/Screenshot_89.png)