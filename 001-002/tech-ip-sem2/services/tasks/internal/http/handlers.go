package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	"tech-ip-sem2/services/tasks/internal/models"
	"tech-ip-sem2/services/tasks/internal/service"
)

type Handlers struct {
	tasksService *service.TasksService
	authClient   *authclient.Client
}

func NewHandlers(tasksService *service.TasksService, authClient *authclient.Client) *Handlers {
	return &Handlers{
		tasksService: tasksService,
		authClient:   authClient,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handlers) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Error: "unauthorized"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Error: "invalid authorization format"})
			return
		}

		token := parts[1]

		valid, subject, err := h.authClient.VerifyToken(r.Context(), token)

		if err != nil {
			errMsg := err.Error()
			switch {
			case strings.Contains(errMsg, "timeout"):
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(errorResponse{Error: "auth service timeout"})
				return
			case strings.Contains(errMsg, "unavailable"):
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadGateway)
				json.NewEncoder(w).Encode(errorResponse{Error: "auth service unavailable"})
				return
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
				return
			}
		}

		if !valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Error: "unauthorized"})
			return
		}

		ctx := context.WithValue(r.Context(), "subject", subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request format"})
		return
	}

	if req.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "title is required"})
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}

	created := h.tasksService.Create(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.tasksService.GetAll()

	type listItem struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}

	// ИСПРАВЛЕНИЕ: создаем срез с длиной 0 и емкостью len(tasks)
	response := make([]listItem, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, listItem{
			ID:    task.ID,
			Title: task.Title,
			Done:  task.Done,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Логируем ошибку, но не отправляем повторно статус (уже отправлен 200)
		// В реальном приложении здесь должен быть логгер
	}
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "task id is required"})
		return
	}

	task, exists := h.tasksService.GetByID(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: "task not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "task id is required"})
		return
	}

	var updates models.TaskUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request format"})
		return
	}

	task, exists := h.tasksService.Update(id, updates)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: "task not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "task id is required"})
		return
	}

	deleted := h.tasksService.Delete(id)
	if !deleted {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: "task not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}