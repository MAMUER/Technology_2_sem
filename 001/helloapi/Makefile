run:
	go run ./cmd/server

build:
	go build -o helloapi.exe ./cmd/server

check:	
	go vet ./...
	go fmt ./...

tree:
	tree /f

fast:
	go vet ./...
	go fmt ./...
	go build -o helloapi.exe ./cmd/server
	go run ./cmd/server