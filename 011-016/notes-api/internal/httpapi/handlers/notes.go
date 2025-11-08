package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/notes-api/internal/core/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service *service.NoteService
}

// Простой middleware
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// CreateNote создает новую заметку
// @Summary Create a new note
// @Description Create a new note with title and content
// @Tags notes
// @Accept json
// @Produce json
// @Param note body CreateNoteRequest true "Note object"
// @Success 201 {object} CreateNoteResponse
// @Router /notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateNote(req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// ListNotes возвращает все заметки
// @Summary Get all notes
// @Description Get all notes
// @Tags notes
// @Produce json
// @Success 200 {array} core.Note
// @Router /notes [get]
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.Service.GetAllNotes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// GetNote возвращает заметку по ID
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := h.Service.GetNoteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// GetNoteByID возвращает заметку по ID
// @Summary Get note by ID
// @Description Get note by ID
// @Tags notes
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} core.Note
// @Router /notes/{id} [get]
func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := h.Service.GetNoteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// UpdateNote обновляет заметку
// @Summary Update note
// @Description Update note by ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Param note body UpdateNoteRequest true "Note object"
// @Success 200 {object} core.Note
// @Router /notes/{id} [put]
func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedNote, err := h.Service.UpdateNote(id, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedNote)
}

// DeleteNote удаляет заметку
// @Summary Delete note
// @Description Delete note by ID
// @Tags notes
// @Param id path int true "Note ID"
// @Success 204 "No Content"
// @Router /notes/{id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteNote(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PatchNote обновляет заметку (частичное обновление)
func (h *Handler) PatchNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title   *string `json:"title,omitempty"`
		Content *string `json:"content,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	currentNote, err := h.Service.GetNoteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	title := currentNote.Title
	if req.Title != nil {
		title = *req.Title
	}

	content := currentNote.Content
	if req.Content != nil {
		content = *req.Content
	}

	updatedNote, err := h.Service.UpdateNote(id, title, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedNote)
}

// ListNotesWithPagination - keyset пагинация
// @Summary Get notes with keyset pagination
// @Description Get notes using keyset pagination for better performance
// @Tags notes
// @Produce json
// @Param cursor_time query string false "Cursor time (RFC3339)"
// @Param cursor_id query int false "Cursor ID"
// @Param limit query int false "Limit (default 20)"
// @Success 200 {array} core.Note
// @Router /notes/paginated [get]
func (h *Handler) ListNotesWithPagination(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var cursorTime time.Time
	var cursorID int64

	if cursorTimeStr := r.URL.Query().Get("cursor_time"); cursorTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, cursorTimeStr); err == nil {
			cursorTime = t
		}
	}

	if cursorIDStr := r.URL.Query().Get("cursor_id"); cursorIDStr != "" {
		if id, err := strconv.ParseInt(cursorIDStr, 10, 64); err == nil {
			cursorID = id
		}
	}

	notes, err := h.Service.GetAllNotesWithPagination(cursorTime, cursorID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// GetNotesBatch - батчинг для получения нескольких заметок
// @Summary Get multiple notes by IDs
// @Description Get multiple notes in single batch request
// @Tags notes
// @Produce json
// @Param ids query string true "Comma-separated list of IDs"
// @Success 200 {array} core.Note
// @Router /notes/batch [get]
func (h *Handler) GetNotesBatch(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "ids parameter is required", http.StatusBadRequest)
		return
	}

	var ids []int64
	for _, idStr := range strings.Split(idsParam, ",") {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		http.Error(w, "no valid IDs provided", http.StatusBadRequest)
		return
	}

	notes, err := h.Service.GetNotesByIDs(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateNoteResponse struct {
	ID int64 `json:"id"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
