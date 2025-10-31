run:
	go run ./cmd/server

build:
	go build -o bin/server.exe ./cmd/server

test:
	go test ./internal/api/...

check:	
	go vet ./...
	go fmt ./...

tree:
ifeq ($(DETECTED_OS),Windows)
	$(TREE_CMD)
else
	@echo "Project structure:"
	@find . -type f -name "*.go" | head -20
	@echo ""
	@echo "Directories:"
	@ls -la
endif

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/server.exe ./cmd/server
	go run ./cmd/server

# Создание .env файла с нужным портом
env:
	@echo "Использование: make env PORT=####"
	@echo PORT=$(PORT) > .env
	@echo "Создан .env файл с PORT=$(PORT)"

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  run       - Запустить сервер (использует .env или порт 8080 по умолчанию)"
	@echo "  build     - Собрать исполняемый файл"
	@echo "  check     - Проверить код и форматирование"
	@echo "  fast      - Проверить, собрать и запустить"
	@echo "  tree      - Показать структуру проекта"
	@echo "  env       - Создать .env файл (использование: make env PORT=8081)"
	@echo "  help      - Показать эту справку"