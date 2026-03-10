package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"tech-ip-sem2/services/tasks/internal/rabbitmq"
	"tech-ip-sem2/shared/logger"
	"tech-ip-sem2/shared/middleware"
)

type JobHandlers struct {
	jobPublisher *rabbitmq.JobPublisher
	log          *logger.Logger
}

type processTaskRequest struct {
	TaskID string `json:"task_id"`
}

type jobResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

func NewJobHandlers(jobPublisher *rabbitmq.JobPublisher, log *logger.Logger) *JobHandlers {
	return &JobHandlers{
		jobPublisher: jobPublisher,
		log:          log,
	}
}

func (h *JobHandlers) ProcessTaskJob(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	log := h.log.WithRequestID(requestID)

	var req processTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request format"})
		return
	}

	if req.TaskID == "" {
		log.Warn("missing task_id in request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "task_id is required"})
		return
	}

	// Генерация уникального message_id
	messageID := uuid.New().String()

	// Публикация задания в очередь
	err := h.jobPublisher.PublishJob(r.Context(), req.TaskID, messageID)
	if err != nil {
		log.Error("failed to publish job", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "failed to queue job"})
		return
	}

	log.Info("job published",
		zap.String("task_id", req.TaskID),
		zap.String("message_id", messageID),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(jobResponse{
		MessageID: messageID,
		Status:    "queued",
	})
}

// Ready – эндпоинт для проверки готовности
func (h *JobHandlers) Ready(w http.ResponseWriter, r *http.Request) {
	if h.jobPublisher != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
	}
}
