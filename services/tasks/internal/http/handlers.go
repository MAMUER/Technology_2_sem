package http

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"tech-ip-sem2/services/tasks/internal/client/authclient"
	"tech-ip-sem2/services/tasks/internal/models"
	"tech-ip-sem2/services/tasks/internal/service"
	"tech-ip-sem2/shared/logger"
	"tech-ip-sem2/shared/middleware"
)

type Handlers struct {
	tasksService *service.TasksService
	authClient   *authclient.Client
	log          *logger.Logger
}

func NewHandlers(tasksService *service.TasksService, authClient *authclient.Client, log *logger.Logger) *Handlers {
	return &Handlers{
		tasksService: tasksService,
		authClient:   authClient,
		log:          log,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handlers) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestID(r.Context())
		log := h.log.WithRequestID(requestID)

		instanceID := os.Getenv("INSTANCE_ID")
		if instanceID == "" {
			instanceID = "unknown"
		}
		w.Header().Set("X-Instance-ID", instanceID)

		var subject string
		var authenticated bool

		// Аутентификация через Bearer token
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token := parts[1]
				valid, subj, err := h.authClient.VerifyToken(r.Context(), token)

				if err == nil && valid {
					subject = subj
					authenticated = true
					log.Info("token authenticated",
						zap.String("subject", subject),
						zap.String("instance", instanceID))
				} else if err != nil {
					log.Error("token verification failed", zap.Error(err))
				}
			}
		}

		// Аутентификация через session cookie
		if !authenticated {
			sessionCookie, err := r.Cookie("session_id")
			if err == nil && sessionCookie.Value != "" {
				subject = "student"
				authenticated = true
				log.Info("cookie authenticated",
					zap.String("subject", subject),
					zap.String("instance", instanceID))
			}
		}

		if !authenticated {
			log.Warn("authentication failed", zap.String("instance", instanceID))
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
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)
	subject := r.Context().Value("subject").(string)

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request format"})
		return
	}

	if req.Title == "" {
		log.Warn("missing title in request")
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

	created, err := h.tasksService.Create(task, subject)
	if err != nil {
		log.Error("failed to create task", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	log.Info("task created",
		zap.String("task_id", created.ID),
		zap.String("title", created.Title),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)
	subject := r.Context().Value("subject").(string)

	tasks, err := h.tasksService.GetAll(subject)
	if err != nil {
		log.Error("failed to get tasks", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	type listItem struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}

	response := make([]listItem, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, listItem{
			ID:    task.ID,
			Title: task.Title,
			Done:  task.Done,
		})
	}

	log.Debug("tasks listed", zap.Int("count", len(response)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("failed to encode response", zap.Error(err))
	}
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)
	subject := r.Context().Value("subject").(string)

	id := r.PathValue("id")
	if id == "" {
		log.Warn("missing task id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "task id is required"})
		return
	}

	task, err := h.tasksService.GetByID(id, subject)
	if err != nil {
		log.Error("failed to get task", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	if task.ID == "" {
		log.Info("task not found", zap.String("task_id", id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: "task not found"})
		return
	}

	log.Debug("task retrieved", zap.String("task_id", id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)
	subject := r.Context().Value("subject").(string)

	id := r.PathValue("id")
	if id == "" {
		log.Warn("missing task id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "task id is required"})
		return
	}

	var updates models.TaskUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Error("failed to decode request", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request format"})
		return
	}

	task, err := h.tasksService.Update(id, updates, subject)
	if err != nil {
		log.Error("failed to update task", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	if task.ID == "" {
		log.Info("task not found for update", zap.String("task_id", id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: "task not found"})
		return
	}

	log.Info("task updated", zap.String("task_id", id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)
	subject := r.Context().Value("subject").(string)

	id := r.PathValue("id")
	if id == "" {
		log.Warn("missing task id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "task id is required"})
		return
	}

	deleted, err := h.tasksService.Delete(id, subject)
	if err != nil {
		log.Error("failed to delete task", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	if !deleted {
		log.Info("task not found for deletion", zap.String("task_id", id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: "task not found"})
		return
	}

	log.Info("task deleted", zap.String("task_id", id))
	w.WriteHeader(http.StatusNoContent)
}

// Поиск задач по названию
func (h *Handlers) SearchTasks(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)
	subject := r.Context().Value("subject").(string)

	term := r.URL.Query().Get("q")
	if term == "" {
		log.Warn("missing search term")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "search term is required"})
		return
	}

	// Уязвимая версия
	useVulnerable := r.URL.Query().Get("vulnerable") == "true"

	var tasks []models.Task
	var err error

	if useVulnerable {
		log.Warn("Using VULNERABLE search - FOR DEMO ONLY", zap.String("term", term))
		tasks, err = h.tasksService.SearchByTitleVulnerable(term, subject)
	} else {
		tasks, err = h.tasksService.SearchByTitle(term, subject)
	}

	if err != nil {
		log.Error("failed to search tasks", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "internal server error"})
		return
	}

	type searchResult struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		Done        bool   `json:"done"`
	}

	response := make([]searchResult, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, searchResult{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Done:        task.Done,
		})
	}

	log.Info("search completed",
		zap.String("term", term),
		zap.Int("results", len(response)),
		zap.Bool("vulnerable", useVulnerable),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Health check
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		instanceID = "unknown"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"service":   "tasks",
		"instance":  instanceID,
		"timestamp": time.Now().Unix(),
	})
}
