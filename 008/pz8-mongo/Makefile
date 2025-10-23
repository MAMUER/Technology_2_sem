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
	tree /f

#ssh root@193.233.175.221
#ssh -L 27017:localhost:27017 root@193.233.175.221 -N
#ssh -L 27017:localhost:27017 root@193.233.175.221 -N -f

#Пользователь MongoDB
#логин: root
#пароль: secret
#порт: 27017
#база данных: pz8

# Остановка туннелей
tunnel-stop:
	taskkill /F /IM ssh.exe 2>nul || echo No SSH tunnels found
	taskkill /F /IM cmd.exe 2>nul || echo No CMD windows found

# Инструкция
setup-teacher:
	@echo === Установка для преподавателя ===
	@echo 1. Запустите SSH тоннель
	@echo 2. Зайдите на сервер
	@echo 3. Проверьте статус MongoDB (При необходимости - запустите MongoDB)
	@echo 4. В другом окне: make run
	@echo 3. Для остановки туннеля: make tunnel-stop

# Статус туннелей
tunnel-status:
	@echo === Active SSH Processes ===
	@tasklist | findstr ssh || echo No active SSH processes found
	@echo === Listening Ports ===
	@netstat -an | findstr 27017 || echo Port 27017 is not listening

.PHONY: run build tunnel tunnel-background tunnel-stop check-db setup-teacher tunnel-status