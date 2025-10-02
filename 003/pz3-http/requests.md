# Тестовые запросы API для задач (Windows PowerShell)

## 1. Проверка здоровья сервера GET /health
Invoke-RestMethod -Uri http://localhost:8080/health -Method GET

## 2. Создать задачу POST /tasks (Не поддерживает русские символы)
Invoke-RestMethod -Uri http://localhost:8080/tasks -Method POST -ContentType "application/json" -Body '{"title":"Новая задача"}'

## 2.1 Создать задачу POST /tasks (Поддерживает русские символы)
$json = '{"title":"Новая задача"}'
$bytes = [System.Text.Encoding]::UTF8.GetBytes($json)
Invoke-RestMethod -Uri http://localhost:8080/tasks -Method POST -ContentType 'application/json' -Body $bytes

## 3. Получить список задач GET /tasks
Invoke-RestMethod -Uri http://localhost:8080/tasks -Method GET

## 4. Получить задачу по ID GET /tasks/1
Invoke-RestMethod -Uri http://localhost:8080/tasks/1 -Method GET

## 5. Отметить задачу выполненной PATCH /tasks/1
Invoke-RestMethod -Uri http://localhost:8080/tasks/1 -Method PATCH

## 6. Поиск задач с фильтром GET /tasks?q=TEXT
Invoke-RestMethod -Uri "http://localhost:8080/tasks?q=TEXT" -Method GET

## 7. Удалить задачу DELETE /tasks/1
Invoke-RestMethod -Uri http://localhost:8080/tasks/1 -Method DELETE

## 8. Создать задачу с пустым заголовком (для проверки ошибки)
Invoke-RestMethod -Uri http://localhost:8080/tasks -Method POST -ContentType "application/json" -Body '{"title":""}'
