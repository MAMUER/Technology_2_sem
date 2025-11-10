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

# SSH туннели
tunnel:
	ssh -L 5433:localhost:5432 root@193.233.175.221 -N -o ServerAliveInterval=30

tunnel-background:
	ssh -L 5433:localhost:5432 root@193.233.175.221 -N -f -o ServerAliveInterval=30

# Остановка туннелей
tunnel-stop:
	$(KILL_SSH)
	$(KILL_CMD)

# Проверка подключения к БД
check-db-root:
	psql postgres://postgres:1@localhost:5433/pz6_gorm?sslmode=disable -c "SELECT version();"

check-db:
	psql postgres://teacher_app:secure_password_123@localhost:5433/pz6_gorm?sslmode=disable -c "SELECT version();"

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@if [ "$(OS)" = "Windows_NT" ]; then \
		tasklist | $(FIND_STR) ssh || echo No active SSH processes found; \
	else \
		ps aux | $(FIND_STR) ssh | $(FIND_STR) -v grep || echo No active SSH processes found; \
	fi
	@echo === Listening Ports ===
	@if [ "$(OS)" = "Windows_NT" ]; then \
		$(NETSTAT) | $(FIND_STR) 5433 || echo Port 5433 is not listening; \
	else \
		$(NETSTAT) || echo Port 5433 is not listening; \
	fi

# Дополнительные команды
clean:
	rm -rf bin/

test:
	go test ./...

help:
	@echo "=== Доступные команды ==="
	@echo "run           - Запустить сервер"
	@echo "build         - Собрать бинарник"
	@echo "fast          - Проверить, собрать и запустить"
	@echo "check         - Проверить код (vet + fmt)"
	@echo "tree          - Показать структуру проекта"
	@echo "tunnel        - Запустить SSH туннель к БД"
	@echo "tunnel-stop   - Остановить SSH туннели"
	@echo "tunnel-status - Показать статус туннелей"
	@echo "check-db      - Проверить подключение к БД"
	@echo "setup-teacher - Инструкция для преподавателя"
	@echo "clean         - Очистить бинарные файлы"
	@echo "test          - Запустить тесты"

.PHONY: run build check tree fast tunnel tunnel-background tunnel-stop check-db setup-teacher tunnel-status clean test help