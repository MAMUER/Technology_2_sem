package main

import (
	"example.com/notes-api/internal/config"
	"example.com/notes-api/internal/core/service"
	"example.com/notes-api/internal/http/handlers"
	"example.com/notes-api/internal/repo"
	"log"
	"net/http"

	httpx "example.com/notes-api/internal/http"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Подключаемся к PostgreSQL через SSH туннель (порт 5433)
	dbPool, err := config.NewDBPool()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	// Создаем репозиторий PostgreSQL
	noteRepo := repo.NewNoteRepoPostgres(dbPool)

	// Создаем сервис
	noteService := service.NewNoteService(noteRepo)

	// Создаем хендлер
	h := &handlers.Handler{Service: noteService}
	r := httpx.NewRouter(h)

	// Swagger документация
	r.Get("/docs/*", httpSwagger.WrapHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
