# Основные команды
run:
	.\bin\server

build:
	go build -o bin/server.exe ./cmd/api

check:	
	go vet ./...
	go fmt ./...

test:
	go test -bench=. -benchmem ./...

tree:
	tree /f