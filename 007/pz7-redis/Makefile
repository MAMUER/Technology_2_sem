# Определение ОС
ifeq ($(OS),Windows_NT)
    # Windows
    BINARY_NAME=server.exe
    RUN_CMD=.\bin\$(BINARY_NAME)
    TREE_CMD=tree /f
    KILL_SSH=taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found
    KILL_CMD=taskkill /F /IM cmd.exe 2>nul || echo No CMD windows found
    FIND_STR=findstr
    NETSTAT=netstat -an
else
    # Linux/Mac
    BINARY_NAME=server
    RUN_CMD=./bin/$(BINARY_NAME)
    TREE_CMD=tree -I 'bin|vendor'
    KILL_SSH=pkill ssh 2>/dev/null || echo No SSH tunnels found
    KILL_CMD=pkill -f "ssh.*6379" 2>/dev/null || echo No SSH processes found
    FIND_STR=grep
    NETSTAT=netstat -an | grep 6379 || ss -tulpn | grep 6379
endif

# Основные команды
run:
	$(RUN_CMD)

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

check:	
	go vet ./...
	go fmt ./...

# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    TREE_CMD := tree /f
else
    DETECTED_OS := Linux
    TREE_CMD := tree
endif

tree:
	@$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/$(BINARY_NAME) ./cmd/server
	$(RUN_CMD)

# SSH туннели
tunnel:
	ssh -L 6379:localhost:6379 root@193.233.175.221 -N -o ServerAliveInterval=30

tunnel-background:
	ssh -L 6379:localhost:6379 root@193.233.175.221 -N -f -o ServerAliveInterval=30

# Остановка туннелей
tunnel-stop:
	$(KILL_SSH)
	$(KILL_CMD)

# Проверка подключения к Redis
check-redis:
	redis-cli -p 6379 ping

check-redis-info:
	redis-cli -p 6379 info server

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
		$(NETSTAT) | $(FIND_STR) 6379 || echo Port 6379 is not listening; \
	else \
		$(NETSTAT) || echo Port 6379 is not listening; \
	fi

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

help:
	@echo "=== Доступные команды (pz7-redis) ==="
	@echo "run           - Запустить сервер"
	@echo "build         - Собрать бинарник"
	@echo "fast          - Проверить, собрать и запустить"
	@echo "check         - Проверить код (vet + fmt)"
	@echo "tree          - Показать структуру проекта"
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

.PHONY: run build check tree fast tunnel tunnel-background tunnel-stop tunnel-status clean test help check-redis check-redis-info redis-flush redis-keys redis-monitor