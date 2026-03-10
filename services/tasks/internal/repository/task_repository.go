package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"tech-ip-sem2/services/tasks/internal/models"
)

type TaskRepository interface {
	Create(task models.Task, subject string) (models.Task, error)
	GetAll(subject string) ([]models.Task, error)
	GetByID(id string, subject string) (models.Task, error)
	Update(id string, updates models.TaskUpdate, subject string) (models.Task, error)
	Delete(id string, subject string) (bool, error)
	SearchByTitle(term string, subject string) ([]models.Task, error)

	// УЯЗВИМАЯ ВЕРСИЯ
	SearchByTitleVulnerable(term string, subject string) ([]models.Task, error)

	Close() error
}

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(connStr string) (*PostgresTaskRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresTaskRepository{
		db: db,
	}, nil
}

func (r *PostgresTaskRepository) Close() error {
	return r.db.Close()
}

// БЕЗОПАСНАЯ ВЕРСИЯ
func (r *PostgresTaskRepository) Create(task models.Task, subject string) (models.Task, error) {
	query := `
        INSERT INTO tasks (id, title, description, due_date, done, subject, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, title, description, due_date, done, subject, created_at, updated_at
    `

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		task.ID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Done,
		subject,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Done,
		&task.Subject,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return models.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (r *PostgresTaskRepository) GetAll(subject string) ([]models.Task, error) {
	query := `
        SELECT id, title, description, due_date, done, subject, created_at, updated_at
        FROM tasks
        WHERE subject = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Done,
			&task.Subject,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(id string, subject string) (models.Task, error) {
	query := `
        SELECT id, title, description, due_date, done, subject, created_at, updated_at
        FROM tasks
        WHERE id = $1 AND subject = $2
    `

	var task models.Task
	err := r.db.QueryRow(query, id, subject).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Done,
		&task.Subject,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.Task{}, nil
	}
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (r *PostgresTaskRepository) Update(id string, updates models.TaskUpdate, subject string) (models.Task, error) {
	task, err := r.GetByID(id, subject)
	if err != nil {
		return models.Task{}, err
	}
	if task.ID == "" {
		return models.Task{}, nil
	}

	if updates.Title != nil {
		task.Title = *updates.Title
	}
	if updates.Description != nil {
		task.Description = *updates.Description
	}
	if updates.DueDate != nil {
		task.DueDate = *updates.DueDate
	}
	if updates.Done != nil {
		task.Done = *updates.Done
	}
	task.UpdatedAt = time.Now()

	query := `
        UPDATE tasks
        SET title = $1, description = $2, due_date = $3, done = $4, updated_at = $5
        WHERE id = $6 AND subject = $7
        RETURNING id, title, description, due_date, done, subject, created_at, updated_at
    `

	err = r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Done,
		task.UpdatedAt,
		id,
		subject,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Done,
		&task.Subject,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return models.Task{}, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (r *PostgresTaskRepository) Delete(id string, subject string) (bool, error) {
	query := `DELETE FROM tasks WHERE id = $1 AND subject = $2`

	result, err := r.db.Exec(query, id, subject)
	if err != nil {
		return false, fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}

// БЕЗОПАСНАЯ ВЕРСИЯ
func (r *PostgresTaskRepository) SearchByTitle(term string, subject string) ([]models.Task, error) {
	query := `
        SELECT id, title, description, due_date, done, subject, created_at, updated_at
        FROM tasks
        WHERE subject = $1 AND title ILIKE $2
        ORDER BY created_at DESC
    `
	rows, err := r.db.Query(query, subject, "%"+term+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Done,
			&task.Subject,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// УЯЗВИМАЯ ВЕРСИЯ
func (r *PostgresTaskRepository) SearchByTitleVulnerable(term string, subject string) ([]models.Task, error) {
	// SQL-инъекция
	query := fmt.Sprintf(`
        SELECT id, title, description, due_date, done, subject, created_at, updated_at
        FROM tasks
        WHERE subject = '%s' AND title ILIKE '%%%s%%'
        ORDER BY created_at DESC
    `, subject, term)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Done,
			&task.Subject,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
