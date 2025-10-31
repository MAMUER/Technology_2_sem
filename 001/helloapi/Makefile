run:
	go run ./cmd/server

build:
	go build -o helloapi.exe ./cmd/server

check:	
	go vet ./...
	go fmt ./...

tree:
	$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o helloapi.exe ./cmd/server
	go run ./cmd/server

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

.PHONY: run build check tree fast env help