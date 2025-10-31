package httpx

import (
	"example.com/notes-api/internal/httpapi/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		// Основные CRUD операции
		r.Get("/notes", h.ListNotes)
		r.Post("/notes", h.CreateNote)
		r.Get("/notes/{id}", h.GetNote)
		r.Patch("/notes/{id}", h.PatchNote)
		r.Put("/notes/{id}", h.UpdateNote)
		r.Delete("/notes/{id}", h.DeleteNote)

		// Новые оптимизированные эндпоинты
		r.Get("/notes/paginated", h.ListNotesWithPagination) // Keyset пагинация
		r.Get("/notes/batch", h.GetNotesBatch)               // Батчинг
	})

	return r
}
