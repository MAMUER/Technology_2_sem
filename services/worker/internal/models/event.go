package models

import "time"

type TaskEvent struct {
	Event     string    `json:"event"`
	TaskID    string    `json:"task_id"`
	Title     string    `json:"title"`
	Subject   string    `json:"subject"`
	Timestamp time.Time `json:"ts"`
	RequestID string    `json:"request_id,omitempty"`
}

const (
	EventTaskCreated = "task.created"
	EventTaskUpdated = "task.updated"
	EventTaskDeleted = "task.deleted"
)
