run:
	.\bin\server

build:
	go build -o bin/server.exe ./

check:	
	go vet ./...
	go fmt ./...

tree:
	tree /f

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/server.exe ./
	.\bin\server