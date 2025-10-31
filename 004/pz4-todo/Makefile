run:
	.\bin\server

build:
	go build -o bin/server.exe ./

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