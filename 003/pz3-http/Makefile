run:
	go run ./cmd/server

build:
	go build -o bin/server.exe ./cmd/server

test:
	go test ./internal/api/...

check:	
	go vet ./...
	go fmt ./...

tree:
	tree /f

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/server.exe ./cmd/server
	go run ./cmd/server