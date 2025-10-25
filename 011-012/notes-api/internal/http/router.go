// internal/http/router.go
package httpx

import (
	"example.com/notes-api/internal/http/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Группа защищенных эндпоинтов
	r.Route("/api/v1", func(r chi.Router) {
		// JWT middleware ко всем эндпоинтам API
		r.Use(h.AuthMiddleware)

		r.Get("/notes", h.ListNotes)
		r.Post("/notes", h.CreateNote)
		r.Get("/notes/{id}", h.GetNote)
		r.Patch("/notes/{id}", h.PatchNote)
		r.Delete("/notes/{id}", h.DeleteNote)
	})

	return r
}
