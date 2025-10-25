# Запуск
run:
	go run main.go

# Бенчмарки
bench:
	go test -bench=. -benchmem ./...

# Профили (запускать при работающем сервере)
profile-cpu:
	go tool pprof -http=:9999 http://localhost:8080/debug/pprof/profile

profile-block:
	go tool pprof -http=:9998 http://localhost:8080/debug/pprof/block

profile-mutex:
	go tool pprof -http=:9997 http://localhost:8080/debug/pprof/mutex

# Нагрузка (нужен hey: go install github.com/rakyll/hey@latest)
load:
	hey -n 100 -c 10 http://localhost:8080/work-slow