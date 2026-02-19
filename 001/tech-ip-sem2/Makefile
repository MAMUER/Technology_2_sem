ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RUN_AUTH_CMD := cd services/auth && set AUTH_PORT=8081&& go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && set TASKS_PORT=8082&& set AUTH_BASE_URL=http://localhost:8081&& go run ./cmd/tasks
    TREE_CMD := tree /f
else
    DETECTED_OS := Linux
    RUN_AUTH_CMD := cd services/auth && AUTH_PORT=8081 go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && TASKS_PORT=8082 AUTH_BASE_URL=http://localhost:8081 go run ./cmd/tasks
    TREE_CMD := tree
endif

AUTH_PORT = 8081
TASKS_PORT = 8082
AUTH_BASE_URL = http://localhost:$(AUTH_PORT)

check:
	go mod tidy
	go vet ./...
	go fmt ./...

fast-auth: check
	@echo Quick start Auth service on port $(AUTH_PORT)...
	$(RUN_AUTH_CMD)

fast-tasks: check
	@echo Quick start Tasks service on port $(TASKS_PORT)...
	$(RUN_TASKS_CMD)

test:
	@echo 1. Get token:
	@curl -s -X POST http://localhost:$(AUTH_PORT)/v1/auth/login -H "Content-Type: application/json" -H "X-Request-ID: test-001" -d "{\"username\":\"student\",\"password\":\"student\"}"
	@echo.
	@echo 2. Verify token:
	@curl -s -X GET http://localhost:$(AUTH_PORT)/v1/auth/verify -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: test-002"
	@echo.
	@echo 3. Create task:
	@curl -s -X POST http://localhost:$(TASKS_PORT)/v1/tasks -H "Content-Type: application/json" -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: test-003" -d "{\"title\":\"Do PZ17\",\"description\":\"split services\",\"due_date\":\"2026-01-10\"}"
	@echo.
	@echo 4. List tasks:
	@curl -s -X GET http://localhost:$(TASKS_PORT)/v1/tasks -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: test-004"

curl-examples:
	@echo 1. Get token:
	@echo curl -X POST http://localhost:$(AUTH_PORT)/v1/auth/login -H "Content-Type: application/json" -H "X-Request-ID: req-001" -d "{\"username\":\"student\",\"password\":\"student\"}"
	@echo.
	@echo 2. Create task:
	@echo curl -X POST http://localhost:$(TASKS_PORT)/v1/tasks -H "Content-Type: application/json" -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: req-002" -d "{\"title\":\"Do PZ17\",\"description\":\"split services\",\"due_date\":\"2026-01-10\"}"
	@echo.
	@echo 3. List tasks:
	@echo curl -X GET http://localhost:$(TASKS_PORT)/v1/tasks -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: req-003"

tree:
	@$(TREE_CMD)

help:
	@echo PORTS:
	@echo   Auth service:    $(AUTH_PORT)
	@echo   Tasks service:   $(TASKS_PORT)
	@echo.
	@echo QUICK RUN:
	@echo   make fast-auth    - Auth with code check
	@echo   make fast-tasks   - Tasks with code check
	@echo.
	@echo TESTING:
	@echo   make test         - Test API
	@echo   make test-cmd     - Show test commands
	@echo   make curl-examples - Show curl examples
	@echo.
	@echo UTILS:
	@echo   make check        - go mod tidy, vet and fmt
	@echo   make tree         - Show project structure
	@echo   make help         - This help

.PHONY: run-auth run-tasks check fast-auth fast-tasks test test-cmd curl-examples tree help