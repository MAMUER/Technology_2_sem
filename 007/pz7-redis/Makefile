# Основные команды
run:
	.\bin\server

build:
	go build -o bin/server.exe ./cmd/server

check:	
	go vet ./...
	go fmt ./...

#ssh root@193.233.175.221
#ssh -L 6379:localhost:6379 root@193.233.175.221 -N -o ServerAliveInterval=30
#ssh -L 6379:localhost:6379 root@193.233.175.221 -N -f -o ServerAliveInterval=30

# Остановка туннелей
tunnel-stop:
	taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found
	taskkill /F /IM cmd.exe 2>nul || echo No CMD windows found

# Инструкция
setup-teacher:
	@echo === Установка для преподавателя ===
	@echo 1. Зайти на сервер, запустить redis
	@echo 2. Запустите SSH тоннель
	@echo 3. В другом окне: make run
	@echo 4. Для остановки туннеля: make tunnel-stop

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@tasklist | findstr ssh || echo No active SSH processes found
	@echo === Listening Ports ===
	@netstat -an | findstr 6379 || echo Port 6379 is not listening

.PHONY: run build tunnel tunnel-background tunnel-stop check-db setup-teacher tunnel-status