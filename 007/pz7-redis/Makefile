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

docker-logs-redis:
	docker-compose logs redis -f

docker-status:
	docker-compose ps

docker-stats:
	docker stats

docker-clean:
	docker-compose down -v
	docker system prune -f

# SSH туннели
tunnel:
	ssh -L 6379:localhost:6379 root@193.233.175.221 -N -o ServerAliveInterval=30

tunnel-background:
	ssh -L 6379:localhost:6379 root@193.233.175.221 -N -f -o ServerAliveInterval=30

# Остановка туннелей
tunnel-stop:
	pkill -f "ssh -L 6379" || true

# Проверка подключения к Redis
check-redis:
	redis-cli -p 6379 ping

check-redis-info:
	redis-cli -p 6379 info server

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@ps aux | grep "ssh -L 6379" | grep -v grep || echo No active SSH processes found
	@echo === Listening Ports ===
	@netstat -tln | grep 6379 || echo Port 6379 is not listening

# Дополнительные команды
clean:
	rm -rf bin/

test:
	go test ./...

# Redis команды
redis-flush:
	redis-cli -p 6379 flushall

redis-keys:
	redis-cli -p 6379 keys "*"

redis-monitor:
	redis-cli -p 6379 monitor

# Redis в Docker
redis-cli:
	docker exec -it pz7-redis redis-cli

help:
	@echo "=== Доступные команды (pz7-redis) ==="
	@echo "run           - Запустить сервер"
	@echo "build         - Собрать бинарник"
	@echo "fast          - Проверить, собрать и запустить"
	@echo "check         - Проверить код (vet + fmt)"
	@echo "tree          - Показать структуру проекта"
	@echo ""
	@echo "=== Docker команды ==="
	@echo "docker-up     - Запустить Docker в фоне"
	@echo "docker-down   - Остановить Docker"
	@echo "docker-logs   - Показать логи всех сервисов"
	@echo "docker-logs-app - Логи приложения"
	@echo "docker-logs-redis - Логи Redis"
	@echo "docker-status - Статус контейнеров"
	@echo "docker-stats  - Статистика ресурсов"
	@echo "docker-clean  - Очистить Docker"
	@echo "redis-cli     - Подключиться к Redis в контейнере"
	@echo ""
	@echo "=== Redis туннели ==="
	@echo "tunnel        - Запустить SSH туннель к Redis"
	@echo "tunnel-stop   - Остановить SSH туннели"
	@echo "tunnel-status - Показать статус туннелей"
	@echo "check-redis   - Проверить подключение к Redis"
	@echo "check-redis-info - Информация о Redis сервере"
	@echo ""
	@echo "=== Redis управление ==="
	@echo "redis-flush   - Очистить все данные Redis"
	@echo "redis-keys    - Показать все ключи Redis"
	@echo "redis-monitor - Мониторинг команд Redis"
	@echo ""
	@echo "=== Утилиты ==="
	@echo "clean         - Очистить бинарные файлы"
	@echo "test          - Запустить тесты"

.PHONY: run build check tree fast \
docker-up docker-down docker-logs docker-logs-app docker-logs-redis docker-status docker-stats docker-clean tunnel \
tunnel-background tunnel-stop tunnel-status clean test help check-redis check-redis-info redis-flush redis-keys redis-monitor redis-cli