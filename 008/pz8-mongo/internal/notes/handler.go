package notes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct{ repo *Repo }

func NewHandler(r *Repo) *Handler { return &Handler{repo: r} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Patch("/{id}", h.patch)
	r.Delete("/{id}", h.del)
	r.Get("/search/text", h.textSearch) // GET /api/v1/notes/search/text?q=query&limit=10
	r.Get("/stats", h.stats)            // GET /api/v1/notes/stats
	r.Get("/stats/daily", h.statsDaily) // GET /api/v1/notes/stats/daily?days=7
	return r
}

func (h *Handler) textSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	skip, _ := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}

	ctx, cancel := reqCtx(r)
	defer cancel()
	notes, err := h.repo.TextSearchWithScore(ctx, q, limit, skip)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"notes": notes,
		"query": q,
		"total": len(notes),
	})
}

func (h *Handler) stats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := reqCtx(r)
	defer cancel()

	stats, err := h.repo.GetStats(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, stats)
}

func (h *Handler) statsDaily(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 {
		days = 7
	} // последние 7 дней

	ctx, cancel := reqCtx(r)
	defer cancel()

	stats, err := h.repo.GetStatsByDay(ctx, days)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"period": days,
		"stats":  stats,
	})
}

func reqCtx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 5*time.Second)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Title == "" {
		writeJSON(w, 400, map[string]string{"error": "invalid_json_or_title"})
		return
	}
	c, cancel := reqCtx(r)
	defer cancel()
	n, err := h.repo.Create(c, in.Title, in.Content)
	if err != nil {
		writeJSON(w, 409, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, n)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, cancel := reqCtx(r)
	defer cancel()
	n, err := h.repo.ByID(c, id)
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, 404, map[string]string{"error": "not_found"})
		return
	}
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, n)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("searchType")
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	skip, _ := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)

	if limit <= 0 || limit > 200 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}

	ctx, cancel := reqCtx(r)
	defer cancel()

	var items interface{}
	var err error

	if searchType == "text" && q != "" {
		items, err = h.repo.TextSearchWithScore(ctx, q, limit, skip)
	} else {
		notes, err := h.repo.List(ctx, q, limit, skip)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		var notesWithScore []NoteWithScore
		for _, note := range notes {
			notesWithScore = append(notesWithScore, NoteWithScore{
				Note:  note,
				Score: 1.0,
			})
		}
		items = notesWithScore
	}

	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"notes":      items,
		"searchType": searchType,
		"query":      q,
	})
}

func (h *Handler) patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	c, cancel := reqCtx(r)
	defer cancel()
	n, err := h.repo.Update(c, id, in.Title, in.Content)
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, 404, map[string]string{"error": "not_found"})
		return
	}
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, n)
}

func (h *Handler) del(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, cancel := reqCtx(r)
	defer cancel()
	if err := h.repo.Delete(c, id); errors.Is(err, ErrNotFound) {
		writeJSON(w, 404, map[string]string{"error": "not_found"})
		return
	} else if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(204)
}
