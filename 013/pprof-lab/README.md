# Практическая работа №13
# Николаенко Михаил ЭФМО-02-21

## Описание проекта и требования

Проект использует инструменты профилирования Go-приложений (pprof), помогает анализировать CPU и memory профили, находить узкие места производительности и проводить оптимизацию на основе метрик.

Для работы необходим Graphiz (версия 14.0.2 и выше) -  для визуализации графов

Для работы с командой make необходимо установить менеджер пакетов Chocolatey и установить команду make

Для работы с командой hey необходимо установить эту команду - для нагрузочного тестирования

Проект на языке Go (необходима версия 1.21 и выше)

## Команды запуска/сборки

### Сборка приложения:

make build

### Запуск приложения:

make run

## Запуск бенчмарков

### Сравнение производительности

make bench

### Профилирование блокировок (через Makefile)

make profile-block

### Профилирование мьютексов (через Makefile)

make profile-mutex

## Команды:

### Нагрузка
hey -n 200 -c 8 http://localhost:8080/work

### Нагрузочное тестирование медленной версии
hey -n 200 -c 8 http://localhost:8080/work-slow

### Нагрузочное тестирование быстрой версии
hey -n 200 -c 8 http://localhost:8080/work-fast

### Тестирование блокировок
hey -n 100 -c 5 http://localhost:8080/block-demo

### Открытие страницы профиля
http://localhost:8080/debug/pprof/

### Скачивание файла profile
http://localhost:8080/debug/pprof/profile?seconds=30 

### Скачивание файла heap
http://localhost:8080/debug/pprof/heap

### Скачивание файла goroutine
http://localhost:8080/debug/pprof/goroutine

### Скачивание файла block
http://localhost:8080/debug/pprof/block

### качивание файла mutex
http://localhost:8080/debug/pprof/mutex

### Анализ CPU профиля
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
#### “топ пожирателей CPU”
- (pprof) top
#### показать исходник с “горячими” строками
- (pprof) list main
#### сгенерировать svg с графом вызовов
- (pprof) web 

### Анализ Block профиля
go tool pprof http://localhost:8080/debug/pprof/block
#### топ операций блокировки
- (pprof) top
#### показать где происходят блокировки
- (pprof) list main   

### Анализ Mutex профиля
go tool pprof http://localhost:8080/debug/pprof/mutex
#### топ конфликтов мьютексов
- (pprof) top
#### показать где происходят блокировки мьютексов
- (pprof) list main   

### CPU профиль с веб-интерфейсом
go tool pprof -http=:9999 http://localhost:8080/debug/pprof/profile?seconds=30

### Heap профиль с веб-интерфейсом
go tool pprof -http=:9998 http://localhost:8080/debug/pprof/heap

### Block профиль с веб-интерфейсом
go tool pprof -http=:9997 http://localhost:8080/debug/pprof/block

### Mutex профиль с веб-интерфейсом
go tool pprof -http=:9996 http://localhost:8080/debug/pprof/mutex

## Структура проекта
```
C:.
│   go.mod
│   Makefile
│   README.md
│
├───bin
│       server.exe
│
├───cmd
│   └───api
│           main.go
│
├───internal
│   └───work
│           slow.go
│           slow_test.go
│           timer.go
│
└───PR13
```

## Скриншоты работы проекта

Инициализация проекта

![фото1](./PR13/Screenshot_1.png)

![фото3](./PR13/Screenshot_3.png)

Проверка и запуск приложения

![фото2](./PR13/Screenshot_2.png)

Нагрузка

![фото4](./PR13/Screenshot_4.png)

![фото5](./PR13/Screenshot_5.png)

![фото23](./PR13/Screenshot_23.png)

Профиль

![фото8](./PR13/Screenshot_8.png)

![фото6](./PR13/Screenshot_6.png)

Скачанные файлы

![фото7](./PR13/Screenshot_7.png)

Отображение нагрузки в Web

![фото9](./PR13/Screenshot_9.png)

![фото10](./PR13/Screenshot_10.png)

![фото13](./PR13/Screenshot_13.png)

![фото12](./PR13/Screenshot_12.png)

![фото14](./PR13/Screenshot_14.png)

![фото18](./PR13/Screenshot_18.png)

![фото19](./PR13/Screenshot_19.png)

![фото22](./PR13/Screenshot_22.png)

Отображение нагрузки в консоли

![фото11](./PR13/Screenshot_11.png)

![фото21](./PR13/Screenshot_21.png)

![фото20](./PR13/Screenshot_20.png)

![фото24](./PR13/Screenshot_24.png)

![фото25](./PR13/Screenshot_25.png)

Test и новый запуск проекта

![фото16](./PR13/Screenshot_16.png)

Структура проекта

![фото15](./PR13/Screenshot_15.png)