// router.go
package httpx

import (
    "example.com/notes-api/internal/httpapi/handlers"
    "net/http/pprof"

    "github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
    r := chi.NewRouter()

    // ✅ Правильное подключение pprof
    r.Route("/debug/pprof", func(r chi.Router) {
        r.Get("/", pprof.Index)
        r.Get("/cmdline", pprof.Cmdline)
        r.Get("/profile", pprof.Profile)
        r.Get("/symbol", pprof.Symbol)
        r.Get("/trace", pprof.Trace)
        r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
        r.Get("/heap", pprof.Handler("heap").ServeHTTP)
        r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
        r.Get("/block", pprof.Handler("block").ServeHTTP)
        r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
    })

    // Остальные роуты...
    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/notes", h.ListNotes)
        r.Post("/notes", h.CreateNote)
        r.Get("/notes/{id}", h.GetNote)
        r.Patch("/notes/{id}", h.PatchNote)
        r.Put("/notes/{id}", h.UpdateNote)
        r.Delete("/notes/{id}", h.DeleteNote)
        r.Get("/notes/paginated", h.ListNotesWithPagination)
        r.Get("/notes/batch", h.GetNotesBatch)
    })

    return r
}