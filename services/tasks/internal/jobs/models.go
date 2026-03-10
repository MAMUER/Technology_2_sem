package jobs
type ProcessTaskJob struct {
	JobType   string `json:"job_type"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
	CreatedAt string `json:"created_at"`
}

const (
	JobTypeProcessTask = "process_task"
	MaxAttempts        = 3
)
