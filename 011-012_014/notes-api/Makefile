# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
    PSQL_CMD := psql
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
    PSQL_CMD := psql
endif

# Основные команды
run: $(BINARY)
	$(RUN_CMD)

build:
	go build -o $(BINARY) ./cmd/api

check:	
	go vet ./...
	go fmt ./...

swagger:
	swag init -g cmd/api/main.go -o docs

tree:
	$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/api
	$(RUN_CMD)

# Docker команды
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v

docker-restart: docker-down docker-up

# Нагрузочное тестирование
load-test: load-test-pagination load-test-batch load-test-single load-test-create

load-test-pagination:
	@echo "=== Testing Pagination (1000 requests, 50 concurrent) ==="
	hey -n 1000 -c 50 "http://localhost:8085/api/v1/notes?limit=20"

load-test-batch:
	@echo "=== Testing Batch Requests (500 requests, 20 concurrent) ==="
	hey -n 500 -c 20 "http://localhost:8085/api/v1/notes?ids=1,2,3,4,5"

load-test-single:
	@echo "=== Testing Single Requests (2000 requests, 100 concurrent) ==="
	hey -n 2000 -c 100 "http://localhost:8085/api/v1/notes/1"

load-test-create:
	@echo "=== Testing Create Operations (300 requests, 10 concurrent) ==="
	hey -n 300 -c 10 -m POST -d "{\"title\":\"Test\",\"content\":\"Content\"}" "http://localhost:8085/api/v1/notes"

# Проверка подключения
check-db:
	$(PSQL_CMD) postgres://teacher_app:secure_password_123@localhost:5434/notes?sslmode=disable -c "SELECT version();"

check-db-docker:
	docker exec postgres-pz11 psql -U teacher_app -d notes -c "SELECT version();"

# Тестовые команды
test-create:
	curl -X POST http://localhost:8085/api/v1/notes \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer test-token" \
	  -d '{"title":"Тестовая заметка","content":"Содержание"}'

test-list:
	curl -X GET "http://localhost:8085/api/v1/notes?page=1&limit=10" \
	  -H "Authorization: Bearer test-token" \
	  -H "Content-Type: application/json"

test-get:
	curl -X GET http://localhost:8085/api/v1/notes/1 \
	  -H "Authorization: Bearer test-token" \
	  -H "Content-Type: application/json"

test-update:
	curl -X PATCH http://localhost:8085/api/v1/notes/1 \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer test-token" \
	  -d '{"title":"Обновленный заголовок"}'

test-delete:
	curl -X DELETE http://localhost:8085/api/v1/notes/1 \
	  -H "Authorization: Bearer test-token" \
	  -H "Content-Type: application/json"

test-swagger:
	curl http://localhost:8085/docs/doc.json

.PHONY: run build check swagger tree fast \
        docker-up docker-down docker-logs docker-clean docker-restart \
        load-test load-test-pagination load-test-batch load-test-single load-test-create \
        check-db check-db-docker test-create test-list test-get test-update test-delete test-swagger