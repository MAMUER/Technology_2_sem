# Detect OS
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
	go build -o $(BINARY) ./

check:	
	go vet ./...
	go fmt ./...

tree:
	@$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./
	$(RUN_CMD)

test:
	go test ./...

# Остановка туннелей
tunnel-stop:
	taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found
	taskkill /F /IM cmd.exe 2>nul || echo No CMD windows found

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@tasklist | findstr ssh || echo No active SSH processes found
	@echo.
	@echo === Listening Ports ===
	@netstat -an | findstr :5433 || echo Port 5433 is not listening

# Очистка билдов
clean:
	if exist bin\server.exe del bin\server.exe
	if exist bin rmdir bin

# Показать помощь
help:
	@echo Доступные команды:
	@echo   run           - Запустить сервер
	@echo   build         - Собрать бинарник
	@echo   check         - Проверить код (vet + fmt)
	@echo   test          - Запустить тесты
	@echo   tree          - Показать структуру проекта
	@echo   fast          - Быстрая сборка и запуск
	@echo   tunnel-stop   - Остановить SSH туннели
	@echo   tunnel-status - Показать статус туннелей
	@echo   clean         - Очистить билды
	@echo   help          - Показать эту справку

.PHONY: run build check test tree fast tunnel-stop tunnel-status clean help