run:
	go run ./cmd/server

build:
	go build -o bin/server.exe ./cmd/server

test:
	go test ./internal/api/...

check:	
	go vet ./...
	go fmt ./...