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

# Пытаемся найти main пакет в типичных местах
ifneq ("$(wildcard ./cmd/myapp/main.go)","")
    MAIN_PACKAGE := ./cmd/myapp
else ifneq ("$(wildcard ./cmd/server/main.go)","")
    MAIN_PACKAGE := ./cmd/server
else ifneq ("$(wildcard ./main.go)","")
    MAIN_PACKAGE := .
else
    MAIN_PACKAGE := .
endif

run:
	$(RUN_CMD)

build:
	go build -o $(BINARY) $(MAIN_PACKAGE)

check:	
	go vet ./...
	go fmt ./...

tree:
	@$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) $(MAIN_PACKAGE)
	$(RUN_CMD)

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
	@echo "  tree      - Показать структуру проекта"
	@echo "  fast      - Проверить, собрать и запустить"
	@echo "  env       - Создать .env файл (использование: make env PORT=8081)"
	@echo "  help      - Показать эту справку"

.PHONY: tree