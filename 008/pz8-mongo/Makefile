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

docker-logs-mongo:
	docker-compose logs mongo -f

docker-status:
	docker-compose ps

docker-stats:
	docker stats

docker-clean:
	docker-compose down -v
	docker system prune -f

# MongoDB команды
mongo-shell:
	docker exec -it mongo-pz8 mongosh -u root -p secret --authenticationDatabase admin pz8

mongo-backup:
	docker exec mongo-pz8 mongodump --uri="mongodb://root:secret@localhost:27017/pz8?authSource=admin" --out=/backup

mongo-restore:
	docker exec mongo-pz8 mongorestore --uri="mongodb://root:secret@localhost:27017/pz8?authSource=admin" /backup/pz8

# Остановка туннелей
tunnel-stop:
	taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found
	taskkill /F /IM cmd.exe 2>nul || echo No CMD windows found

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@tasklist | findstr ssh || echo No active SSH processes found
	@echo === Listening Ports ===
	@netstat -an | findstr 27017 || echo Port 27017 is not listening

# Дополнительные команды
clean:
	rm -rf bin/

test:
	go test ./...

help:
	@echo "=== Доступные команды (pz8-mongo) ==="
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
	@echo "docker-logs-mongo - Логи MongoDB"
	@echo "docker-status - Статус контейнеров"
	@echo "docker-stats  - Статистика ресурсов"
	@echo "docker-clean  - Очистить Docker"
	@echo ""
	@echo "=== MongoDB команды ==="
	@echo "mongo-shell   - Подключиться к MongoDB shell"
	@echo "mongo-backup  - Создать бэкап базы данных"
	@echo "mongo-restore - Восстановить базу данных"
	@echo ""
	@echo "=== Утилиты ==="
	@echo "clean         - Очистить бинарные файлы"
	@echo "test          - Запустить тесты"

.PHONY: run build check tree fast \
docker-up docker-down docker-logs docker-logs-app docker-logs-mongo docker-status docker-stats docker-clean \
tunnel-stop tunnel-status clean test help mongo-shell mongo-backup mongo-restore