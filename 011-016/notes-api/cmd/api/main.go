package main

import (
	"example.com/notes-api/internal/config"
	"example.com/notes-api/internal/core/service"
	"example.com/notes-api/internal/httpapi/handlers"
	"example.com/notes-api/internal/repo"
	"example.com/notes-api/internal/work"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"

	httpx "example.com/notes-api/internal/httpapi"
	httpSwagger "github.com/swaggo/http-swagger"

	"example.com/notes-api/internal/mathx"
	userservice "example.com/notes-api/internal/service"
	"example.com/notes-api/internal/stringsx"
)

var (
	globalMutex sync.Mutex
	sharedData  int
)

// Простая реализация репозитория для демонстрации
type demoRepo struct{}

func (d demoRepo) ByEmail(email string) (userservice.UserRecord, error) {
	// Демо данные
	users := map[string]userservice.UserRecord{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: "demo_hash"},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: "demo_hash"},
	}

	user, exists := users[email]
	if !exists {
		return userservice.UserRecord{}, userservice.ErrNotFound
	}
	return user, nil
}

func main() {
	// Включаем профилирование блокировок и мьютексов
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	log.Println("Block and mutex profiling enabled")

	// Подключаемся к PostgreSQL
	dbPool, err := config.NewDBPool()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	// Создаем репозиторий PostgreSQL
	noteRepo := repo.NewNoteRepoPostgres(dbPool)

	// Создаем сервис заметок
	noteService := service.NewNoteService(noteRepo)

	// Создаем хендлер заметок
	h := &handlers.Handler{Service: noteService}
	r := httpx.NewRouter(h)

	// Swagger документация
	r.Get("/docs/*", httpSwagger.WrapHandler)

	// 🔥 PPROF эндпоинты

	// Медленная версия Fibonacci
	r.HandleFunc("/work-slow", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("Fib(38)")()
		res := work.Fib(38)
		fmt.Fprintf(w, "Slow: %d\n", res)
	})

	// Быстрая версия Fibonacci
	r.HandleFunc("/work-fast", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("FibFast(38)")()
		res := work.FibFast(38)
		fmt.Fprintf(w, "Fast: %d\n", res)
	})

	// Демо блокировок
	r.HandleFunc("/block-demo", func(w http.ResponseWriter, r *http.Request) {
		// Блокировка каналом
		ch := make(chan int)
		go func() {
			ch <- 42
		}()
		result := <-ch

		// Конкуренция за мьютекс
		var wg sync.WaitGroup
		for range 50 {
			wg.Go(func() {
				globalMutex.Lock()
				sharedData++
				globalMutex.Unlock()
			})
		}
		wg.Wait()

		fmt.Fprintf(w, "Result: %d, Shared: %d\n", result, sharedData)
	})

	// 🔥 Эндпоинты из второго проекта

	// Math Operations
	r.HandleFunc("/math/sum/{a}/{b}", func(w http.ResponseWriter, r *http.Request) {
		a := r.PathValue("a")
		b := r.PathValue("b")
		var numA, numB int
		fmt.Sscanf(a, "%d", &numA)
		fmt.Sscanf(b, "%d", &numB)

		sum := mathx.Sum(numA, numB)
		fmt.Fprintf(w, "Sum(%d, %d) = %d\n", numA, numB, sum)
	})

	r.HandleFunc("/math/divide/{a}/{b}", func(w http.ResponseWriter, r *http.Request) {
		a := r.PathValue("a")
		b := r.PathValue("b")
		var numA, numB int
		fmt.Sscanf(a, "%d", &numA)
		fmt.Sscanf(b, "%d", &numB)

		result, err := mathx.Divide(numA, numB)
		if err != nil {
			fmt.Fprintf(w, "Divide(%d, %d) error: %v\n", numA, numB, err)
		} else {
			fmt.Fprintf(w, "Divide(%d, %d) = %d\n", numA, numB, result)
		}
	})

	// String Operations
	r.HandleFunc("/strings/clip/{text}/{max}", func(w http.ResponseWriter, r *http.Request) {
		text := r.PathValue("text")
		max := r.PathValue("max")
		var maxLen int
		fmt.Sscanf(max, "%d", &maxLen)

		clipped := stringsx.Clip(text, maxLen)
		fmt.Fprintf(w, "Clip('%s', %d) = '%s'\n", text, maxLen, clipped)
	})

	// User Service
	r.HandleFunc("/users/find/{email}", func(w http.ResponseWriter, r *http.Request) {
		email := r.PathValue("email")

		repo := demoRepo{}
		svc := userservice.New(repo)

		id, err := svc.FindIDByEmail(email)
		if err != nil {
			fmt.Fprintf(w, "User '%s' not found: %v\n", email, err)
		} else {
			fmt.Fprintf(w, "Found user '%s' with ID: %d\n", email, id)
		}
	})

	// Демо страница со всеми операциями
	r.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `Available Demo Endpoints:
		
Math Operations:
  GET /math/sum/{a}/{b}           - Sum two numbers
  GET /math/divide/{a}/{b}        - Divide two numbers

String Operations:
  GET /strings/clip/{text}/{max}  - Clip text to max length

User Operations:
  GET /users/find/{email}         - Find user by email

Performance Tests:
  GET /work-slow                  - Slow Fibonacci (recursive)
  GET /work-fast                  - Fast Fibonacci (optimized)
  GET /block-demo                 - Blocking operations demo

Notes API:
  GET /api/v1/notes/*             - CRUD operations for notes

Profiling:
  GET /debug/pprof/               - PPROF main page
  GET /docs/                      - Swagger documentation
`)
	})

	log.Println("Server started at :8080")
	log.Println("Available endpoints:")
	log.Println("  - Notes API: /api/v1/notes/*")
	log.Println("  - Math Operations: /math/*")
	log.Println("  - String Operations: /strings/*")
	log.Println("  - User Service: /users/*")
	log.Println("  - Demo Page: /demo")
	log.Println("  - Swagger: /docs/")
	log.Println("  - PPROF: /debug/pprof/")
	log.Println("  - Performance tests: /work-slow, /work-fast, /block-demo")

	log.Fatal(http.ListenAndServe(":8080", r))
}
