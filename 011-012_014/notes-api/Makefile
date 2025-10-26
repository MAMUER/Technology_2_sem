# Основные команды
run:
	.\bin\server

build:
	go build -o bin/server.exe ./cmd/api

check:	
	go vet ./...
	go fmt ./...

swagger:
	swag init -g cmd/api/main.go -o docs

tree:
	tree /f

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/server.exe ./cmd/api
	.\bin\server

# ssh root@193.233.175.221
# ssh -L 6379:localhost:6379 root@193.233.175.221 -N -o ServerAliveInterval=30
# ssh -L 6379:localhost:6379 root@193.233.175.221 -N -f -o ServerAliveInterval=30

# Нагрузочное тестирование
load-test: load-test-pagination load-test-batch load-test-single load-test-create

load-test-pagination:
	@echo "=== Testing Pagination (1000 requests, 50 concurrent) ==="
	hey -n 1000 -c 50 "http://localhost:8080/api/v1/notes/paginated?limit=20"

load-test-batch:
	@echo "=== Testing Batch Requests (500 requests, 20 concurrent) ==="
	hey -n 500 -c 20 "http://localhost:8080/api/v1/notes/batch?ids=1,2,3,4,5"

load-test-single:
	@echo "=== Testing Single Requests (2000 requests, 100 concurrent) ==="
	hey -n 2000 -c 100 "http://localhost:8080/api/v1/notes/1"

load-test-create:
	@echo "=== Testing Create Operations (300 requests, 10 concurrent) ==="
	hey -n 300 -c 10 -m POST -d "{\"title\":\"Test\",\"content\":\"Content\"}" "http://localhost:8080/api/v1/notes"

# Тюнинг пула соединений
pool-test-small:
	@echo "=== Testing with small pool (MaxConns=5) ==="
	@set DATABASE_POOL_MAX_CONNS=5&& set DATABASE_POOL_MIN_CONNS=2&& go run ./cmd/api

pool-test-medium:
	@echo "=== Testing with medium pool (MaxConns=20) ==="
	@set DATABASE_POOL_MAX_CONNS=20&& set DATABASE_POOL_MIN_CONNS=5&& go run ./cmd/api

pool-test-large:
	@echo "=== Testing with large pool (MaxConns=50) ==="
	@set DATABASE_POOL_MAX_CONNS=50&& set DATABASE_POOL_MIN_CONNS=10&& go run ./cmd/api

# Анализ производительности БД
db-stats:
	@echo "=== Top 10 Slowest Queries ==="
	psql postgres://teacher_app:secure_password_123@localhost:5433/notes?sslmode=disable -c "SELECT query, calls, round(total_exec_time::numeric, 2) as total_time, round(mean_exec_time::numeric, 2) as mean_time, rows FROM pg_stat_statements ORDER BY total_exec_time DESC LIMIT 10;"

db-explain-pagination:
	@echo "=== EXPLAIN Pagination Query ==="
	psql postgres://teacher_app:secure_password_123@localhost:5433/notes?sslmode=disable -c "EXPLAIN (ANALYZE, BUFFERS) SELECT id, title, content, created_at FROM notes WHERE (created_at, id) < (now(), 1000) ORDER BY created_at DESC, id DESC LIMIT 20;"

tunnel-stop:
	taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found

# Проверка подключения
check-db:
	psql postgres://teacher_app:secure_password_123@localhost:5433/notes?sslmode=disable -c "SELECT version();"

# Мониторинг БД в реальном времени
monitor-db:
	@echo "=== Database Connections Monitor (Ctrl+C to stop) ==="
	@while true; do \
		psql postgres://teacher_app:secure_password_123@localhost:5433/notes?sslmode=disable -c "SELECT datname as database, numbackends as connections, xact_commit as commits, xact_rollback as rollbacks FROM pg_stat_database WHERE datname = 'notes';" 2>/dev/null || echo "Database not available"; \
		sleep 5; \
	done

# Очистка
clean:
	@echo "=== Cleaning ==="
	go clean
	rm -f bin/server.exe

.PHONY: run build check swagger tree \
		load-test load-test-pagination load-test-batch load-test-single load-test-create \
		pool-test-small pool-test-medium pool-test-large \
		db-stats db-explain-pagination \
		tunnel tunnel-stop check-db monitor-db clean