run:
	go run ./cmd/myapp

build:
	go build -o bin/myapp.exe ./cmd/myapp

check:	
	go vet ./...
	go fmt ./...

tree:
	tree /f

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/myapp.exe ./cmd/myapp
	go run ./cmd/myapp

# Создание .env файла с нужным портом
env:
	@echo "Usage: make env PORT=####"
	@echo APP_PORT=$(PORT) > .env
	@echo "Created .env with PORT=$(PORT)"

# Помощь
help:
	@echo "Available commands:"
	@echo "  run       - Run server (uses .env or default port 8080)"
	@echo "  build     - Build executable"
	@echo "  check     - Run vet and fmt"
	@echo "  fast      - Check, build and run"
	@echo "  tree      - Show project structure"
	@echo "  env       - Create .env file (usage: make env PORT=8081)"
	@echo "  help      - Show this help"