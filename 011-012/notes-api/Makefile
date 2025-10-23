# Основные команды
run:
	.\bin\server

build:
	go build -o bin/server.exe ./cmd/api

check:	
	go vet ./...
	go fmt ./...

swagger:
	swag init -g cmd/api/main.go -o docs

tree:
	tree /f

.PHONY: run swagger