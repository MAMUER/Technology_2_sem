run:
	.\bin\server

build:
	go build -o bin/server.exe ./

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
	go build -o bin/server.exe ./
	.\bin\server

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  run       - Запустить сервер (использует .env или порт 8080 по умолчанию)"
	@echo "  build     - Собрать исполняемый файл"
	@echo "  check     - Проверить код и форматирование"
	@echo "  fast      - Проверить, собрать и запустить"
	@echo "  tree      - Показать структуру проекта"
	@echo "  help      - Показать эту справку"

.PHONY: tree