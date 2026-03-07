package models

import (
	"tech-ip-sem2/shared/sanitize"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	Done        bool      `json:"done"`
	Subject     string    `json:"-"` // Не возвращаем клиенту
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type TaskUpdate struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	Done        *bool   `json:"done,omitempty"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type SearchTaskRequest struct {
	Term string `json:"term"`
}

// Sanitize очищает поля задачи от потенциально опасного содержимого
func (t *Task) Sanitize() {
	t.Title = sanitize.SanitizeText(t.Title)
	t.Description = sanitize.SanitizeHTML(t.Description)
}

// Validate проверяет корректность задачи
func (t *CreateTaskRequest) Validate() error {
	if t.Title == "" {
		return &ValidationError{"title is required"}
	}

	if len(t.Title) > 255 {
		return &ValidationError{"title too long (max 255 characters)"}
	}

	if len(t.Description) > 1000 {
		return &ValidationError{"description too long (max 1000 characters)"}
	}

	// Очищаем поля
	t.Title = sanitize.SanitizeText(t.Title)
	t.Description = sanitize.SanitizeHTML(t.Description)

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
