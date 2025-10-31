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
	@if command -v tree >/dev/null 2>&1; then \
		$(TREE_CMD); \
	else \
		echo "Tree command not found. Installing..."; \
		if command -v apt-get >/dev/null 2>&1; then \
			sudo apt-get update && sudo apt-get install -y tree; \
		elif command -v yum >/dev/null 2>&1; then \
			sudo yum install -y tree; \
		elif command -v apk >/dev/null 2>&1; then \
			sudo apk add tree; \
		else \
			echo "Cannot install tree automatically. Please install manually: sudo apt-get install tree"; \
			echo "Alternative: using find command:"; \
			find . -type f -name "*.go" | head -20; \
		fi; \
	fi

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