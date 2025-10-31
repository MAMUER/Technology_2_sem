# Основные команды
run:
	.\bin\server

build:
	go build -o bin/server.exe ./cmd/api

check:	
	go vet ./...
	go fmt ./...

test:
	go test ./...

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
	go build -o bin/server.exe ./cmd/api
	.\bin\server

# Остановка туннелей
tunnel-stop:
	taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found
	taskkill /F /IM cmd.exe 2>nul || echo No CMD windows found

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@tasklist | findstr ssh || echo No active SSH processes found
	@echo === Listening Ports ===
	@netstat -an | findstr 27017 || echo Port 27017 is not listening

.PHONY: run build tunnel tunnel-background tunnel-stop check-db setup-teacher tunnel-status