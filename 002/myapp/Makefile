run:
	go run ./cmd/myapp

build:
	go build -o bin/myapp.exe ./cmd/myapp

check:	
	go vet ./...
	go fmt ./...

tree:
	@if command -v tree >/dev/null 2>&1; then \
		$(TREE_CMD); \
	else \
		echo "Tree command not found. Installing..."; \
		if command -v apt-get >/dev/null 2>&1; then \
			sudo apt-get update && sudo apt-get install -y tree; \
		elif command -v yum >/dev/null 2>&1; then \
			sudo yum install -y tree; \
		elif command -v apk >/dev/null 2>&1; then \
			sudo apk add tree; \
		else \
			echo "Cannot install tree automatically. Please install manually: sudo apt-get install tree"; \
			echo "Alternative: using find command:"; \
			find . -type f -name "*.go" | head -20; \
		fi; \
	fi

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/myapp.exe ./cmd/myapp
	go run ./cmd/myapp

# Создание .env файла с нужным портом
env:
	@echo "Использование: make env PORT=####"
	@echo APP_PORT=$(PORT) > .env
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