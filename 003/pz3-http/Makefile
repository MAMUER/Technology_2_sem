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

.PHONY: tree