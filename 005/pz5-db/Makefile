# Основные команды
run:
	.\bin\server

build:
	go build -o bin/server.exe ./

check:	
	go vet ./...
	go fmt ./...

tree:
	tree /f

#ssh root@193.233.175.221
#ssh -L 5433:localhost:5432 root@193.233.175.221 -N -o ServerAliveInterval=30
#ssh -L 5433:localhost:5432 root@193.233.175.221 -N -f -o ServerAliveInterval=30

#Пользователь PostgreSQL
#логин: postgres
#пароль: 1

# Остановка туннелей
tunnel-stop:
	taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found
	taskkill /F /IM cmd.exe 2>nul || echo No CMD windows found

# Проверка подключения
check-db-root:
	psql postgres://postgres:1@localhost:5432/todo?sslmode=disable -c "SELECT version();"

check-db:
	psql postgres://teacher_app:secure_password_123@localhost:5432/todo?sslmode=disable -c "SELECT version();"

# Инструкция
setup-teacher:
	@echo === Установка для преподавателя ===
	@echo 1. Запустите SSH тоннель
	@echo 2. В другом окне: make run
	@echo 3. Для остановки туннеля: make tunnel-stop

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@tasklist | findstr ssh || echo No active SSH processes found
	@echo === Listening Ports ===
	@netstat -an | findstr 5433 || echo Port 5433 is not listening

.PHONY: run build tunnel tunnel-background tunnel-stop check-db setup-teacher tunnel-status