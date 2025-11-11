# Определить ОС
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
    FIND_PROCESS := tasklist | findstr
    FIND_PORT := netstat -an | findstr
    KILL_CMD := taskkill /F /IM
    PSQL_CMD := psql
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
    FIND_PROCESS := ps aux | grep
    FIND_PORT := netstat -tln | grep
    KILL_CMD := pkill -f
    PSQL_CMD := psql
endif

run:
	$(RUN_CMD)

build:
	go build -o $(BINARY) ./cmd/api

check:	
	go vet ./...
	go fmt ./...

tree:
	@$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/api
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

docker-logs-postgres:
	docker-compose logs postgres -f

docker-status:
	docker-compose ps

docker-stats:
	docker stats

docker-clean:
	docker-compose down -v
	docker system prune -f

docker-restart:
	docker-compose restart

# PostgreSQL команды
postgres-shell:
	docker exec -it postgres-pz9 psql -U teacher_app -d todo

postgres-backup:
	docker exec postgres-pz9 pg_dump -U teacher_app todo > backup.sql

postgres-restore:
	docker exec -i postgres-pz9 psql -U teacher_app -d todo < backup.sql

# Остановка туннелей
tunnel-stop:
	$(KILL_CMD) "ssh -L 5433" 2>/dev/null || echo No SSH tunnels found

# Проверка подключения
check-db-root:
	$(PSQL_CMD) postgres://postgres:1@localhost:5433/todo?sslmode=disable -c "SELECT version();"

check-db:
	$(PSQL_CMD) postgres://teacher_app:secure_password_123@localhost:5433/todo?sslmode=disable -c "SELECT version();"

# Проверка подключения к контейнеру
check-db-docker:
	docker exec postgres-pz9 psql -U teacher_app -d todo -c "SELECT version();"

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@$(FIND_PROCESS) "ssh -L 5433" | grep -v grep || echo No active SSH processes found
	@echo === Listening Ports ===
	@$(FIND_PORT) 5433 || echo Port 5433 is not listening
	@$(FIND_PORT) 8083 || echo Port 8083 is not listening

# Дополнительные команды
clean:
	rm -rf bin/

test:
	go test ./...

# Разработка
dev: check build
	$(RUN_CMD)

# Миграции (если будут нужны)
migrate: build
	$(RUN_CMD) -migrate

help:
	@echo "=== Доступные команды (pz9-auth) ==="
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
	@echo "docker-logs-postgres - Логи PostgreSQL"
	@echo "docker-status - Статус контейнеров"
	@echo "docker-stats  - Статистика ресурсов"
	@echo "docker-clean  - Очистить Docker"
	@echo "docker-restart - Перезапустить контейнеры"
	@echo ""
	@echo "=== PostgreSQL команды ==="
	@echo "postgres-shell - Подключиться к PostgreSQL"
	@echo "postgres-backup - Создать бэкап базы"
	@echo "postgres-restore - Восстановить базу"
	@echo ""
	@echo "=== Проверки ==="
	@echo "check-db      - Проверить подключение к БД"
	@echo "check-db-docker - Проверить БД в контейнере"
	@echo ""
	@echo "=== Утилиты ==="
	@echo "clean         - Очистить бинарные файлы"
	@echo "test          - Запустить тесты"
	@echo "dev           - Проверить и запустить"

.PHONY: run build check tree fast \
docker-up docker-down docker-logs docker-logs-app docker-logs-postgres docker-status docker-stats docker-clean docker-restart \
tunnel-stop check-db check-db-docker tunnel-status clean test dev migrate postgres-shell postgres-backup postgres-restore help