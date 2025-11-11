# Определить ОС
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
endif

run:
	$(RUN_CMD)

build:
	go build -o $(BINARY) ./cmd/server

check:	
	go vet ./...
	go fmt ./...

tree:
	@$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/server
	$(RUN_CMD)

# Docker команды
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-logs-app:
	docker-compose logs app -f

docker-status:
	docker-compose ps

docker-stats:
	docker stats

docker-clean:
	docker-compose down -v
	docker system prune -f

docker-restart:
	docker-compose restart

# Мониторинг
monitor:
	watch -n 2 'curl -s http://localhost:8084/health && echo "" && docker stats --no-stream'

# Тестовые команды
test-login-admin:
	curl -X POST http://localhost:8084/api/v1/login \
	  -H "Content-Type: application/json" \
	  -d '{"email":"admin@example.com","password":"secret123"}'

test-login-user:
	curl -X POST http://localhost:8084/api/v1/login \
	  -H "Content-Type: application/json" \
	  -d '{"email":"user@example.com","password":"secret123"}'

test-login-user2:
	curl -X POST http://localhost:8084/api/v1/login \
	  -H "Content-Type: application/json" \
	  -d '{"email":"user2@example.com","password":"secret123"}'

test-refresh:
	@echo "Use refresh token from login response:"
	@echo "curl -X POST http://localhost:8084/api/v1/refresh -H \"Authorization: Bearer <refresh_token>\""

test-protected:
	@echo "Use access token from login response:"
	@echo "curl http://localhost:8084/api/v1/protected -H \"Authorization: Bearer <access_token>\""

test-health:
	curl http://localhost:8084/health

test-all: test-health test-login-admin test-login-user

# Утилиты
clean:
	rm -rf bin/

test:
	go test ./...

help:
	@echo "=== Доступные команды (pz10-auth) ==="
	@echo "run           - Запустить сервер"
	@echo "build         - Собрать бинарник"
	@echo "fast          - Проверить, собрать и запустить"
	@echo "check         - Проверить код (vet + fmt)"
	@echo "tree          - Показать структуру проекта"
	@echo ""
	@echo "=== Docker команды ==="
	@echo "docker-up     - Запустить Docker в фоне"
	@echo "docker-down   - Остановить Docker"
	@echo "docker-logs   - Показать логи"
	@echo "docker-logs-app - Логи приложения"
	@echo "docker-status - Статус контейнеров"
	@echo "docker-stats  - Статистика ресурсов"
	@echo "docker-clean  - Очистить Docker"
	@echo "docker-restart - Перезапустить контейнер"
	@echo ""
	@echo "=== Тестовые запросы ==="
	@echo "test-login-admin - Логин администратора"
	@echo "test-login-user  - Логин пользователя"
	@echo "test-login-user2 - Логин второго пользователя"
	@echo "test-refresh     - Обновить токен (инструкция)"
	@echo "test-protected   - Доступ к защищенному роуту"
	@echo "test-health      - Проверить health endpoint"
	@echo "test-all         - Выполнить все базовые тесты"
	@echo ""
	@echo "=== Утилиты ==="
	@echo "clean         - Очистить бинарные файлы"
	@echo "test          - Запустить тесты"
	@echo "monitor       - Мониторинг в реальном времени"
	@echo ""
	@echo "=== Тестовые пользователи ==="
	@echo "  admin@example.com / secret123 (role: admin)"
	@echo "  user@example.com  / secret123 (role: user)"
	@echo "  user2@example.com / secret123 (role: user)"

.PHONY: run build check tree fast \
docker-up docker-down docker-logs docker-logs-app docker-status docker-stats docker-clean docker-restart \
monitor test-login-admin test-login-user test-login-user2 test-refresh test-protected test-health test-all clean test help