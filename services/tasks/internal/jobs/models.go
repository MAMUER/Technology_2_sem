package jobs

// Убираем неиспользуемый импорт "time"
type ProcessTaskJob struct {
	JobType   string `json:"job_type"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
	CreatedAt string `json:"created_at"` // time теперь не нужен, оставляем string
}

const (
	JobTypeProcessTask = "process_task"
	MaxAttempts        = 3
)
