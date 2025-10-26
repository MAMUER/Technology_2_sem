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