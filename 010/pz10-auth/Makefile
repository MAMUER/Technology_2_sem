# Основные команды
run:
	.\bin\server

build:
	go build -o bin/server.exe ./cmd/server

check:	
	go vet ./...
	go fmt ./...

help:
	@echo "Test users:"
	@echo "  admin@example.com / secret123 (role: admin)"
	@echo "  user@example.com  / secret123 (role: user)"
	@echo "  user2@example.com / secret123 (role: user)"

tree:
	tree /f

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/server.exe ./cmd/server
	.\bin\server