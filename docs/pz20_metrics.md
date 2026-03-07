# Метрики
- http_requests_total - счётчик всех запросов
- labels: method, route, status, service
- http_request_duration_seconds - гистограмма длительности запросов
- labels: method, route, service
- buckets: 0.01, 0.05, 0.1, 0.3, 1, 3, 5, 10 секунд
- http_requests_in_flight - текущее количество активных запросов
- labels: service

# Дашборды в Grafana (4 основных графика)
- RPS (Requests Per Second) по методам
- Количество ошибок 4xx и 5xx
- 95-й перцентиль длительности запросов
- Текущее количество активных запросов

# Prometheus http://193.233.175.221:9090

## RPS (Requests Per Second)

rate(http_requests_total[1m])

## Количество ошибок 4xx

sum(rate(http_requests_total{service="tasks", status=~"4[0-9][0-9]"}[1m])) by (status)

## Количество ошибок 5xx

sum(rate(http_requests_total{service="tasks", status=~"5.."}[1m]))

## 95-й перцентиль длительности запросов

histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{service="tasks"}[5m])) by (le, method))

## Текущие активные запросы

http_requests_in_flight{service="tasks"}

# Grafana http://193.233.175.221:3000
## Доступные дашборды
- Dashboards -> Manage
- "Tasks Service Dashboard"

## Графики в дашборде
- RPS (Requests per second) - по методам GET/POST
- Error Rate - отдельно 4xx и 5xx ошибки
- Request Duration p95 - 95-й перцентиль длительности
- In-Flight Requests - текущие активные запросы

## GET http://193.233.175.221:8082/metrics
- Получение метрик tasks сервиса

Ответ 200:
```
http_requests_in_flight
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",route="/v1/tasks",service="tasks",status="200"} 42
http_requests_total{method="GET",route="/v1/tasks",service="tasks",status="401"} 18
http_requests_total{method="POST",route="/v1/tasks",service="tasks",status="201"} 10

# HELP http_request_duration_seconds Duration of HTTP requests in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.01",method="GET",route="/v1/tasks",service="tasks"} 30
http_request_duration_seconds_bucket{le="0.05",method="GET",route="/v1/tasks",service="tasks"} 45
http_request_duration_seconds_bucket{le="0.1",method="GET",route="/v1/tasks",service="tasks"} 48
http_request_duration_seconds_bucket{le="0.3",method="GET",route="/v1/tasks",service="tasks"} 50
http_request_duration_seconds_bucket{le="+Inf",method="GET",route="/v1/tasks",service="tasks"} 50

# HELP http_requests_in_flight Current number of in-flight HTTP requests
# TYPE http_requests_in_flight gauge
http_requests_in_flight{service="tasks"} 3
```
## GET http://193.233.175.221:8081/metrics
- Получение метрик auth сервиса

Ответ 200:
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="POST",route="/v1/auth/login",service="auth",status="200"} 15
http_requests_total{method="POST",route="/v1/auth/login",service="auth",status="401"} 5
http_requests_total{method="GET",route="/v1/auth/verify",service="auth",status="200"} 25
http_requests_total{method="GET",route="/v1/auth/verify",service="auth",status="401"} 8

# HELP http_request_duration_seconds Duration of HTTP requests in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.01",method="POST",route="/v1/auth/login",service="auth"} 10
http_request_duration_seconds_bucket{le="0.05",method="POST",route="/v1/auth/login",service="auth"} 18
http_request_duration_seconds_bucket{le="0.1",method="POST",route="/v1/auth/login",service="auth"} 20
http_request_duration_seconds_bucket{le="+Inf",method="POST",route="/v1/auth/login",service="auth"} 20

# HELP http_requests_in_flight Current number of in-flight HTTP requests
# TYPE http_requests_in_flight gauge
http_requests_in_flight{service="auth"} 2
```