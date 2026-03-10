package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"tech-ip-sem2/services/graphql/graph/model"
)

type TaskRepository interface {
	Create(task *model.Task, subject string) (*model.Task, error)
	GetAll(subject string) ([]*model.Task, error)
	GetByID(id string, subject string) (*model.Task, error)
	Update(id string, input model.UpdateTaskInput, subject string) (*model.Task, error)
	Delete(id string, subject string) (bool, error)
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

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id VARCHAR(50) PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            due_date DATE,
            done BOOLEAN DEFAULT FALSE,
            subject VARCHAR(100) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresTaskRepository{
		db: db,
	}, nil
}

func (r *PostgresTaskRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresTaskRepository) Create(task *model.Task, subject string) (*model.Task, error) {
	query := `
        INSERT INTO tasks (id, title, description, due_date, done, subject, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, title, description, due_date, done, created_at, updated_at
    `

	now := time.Now()
	createdAt := now.Format(time.RFC3339)
	updatedAt := now.Format(time.RFC3339)

	err := r.db.QueryRow(
		query,
		task.ID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Done,
		subject,
		now,
		now,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Done,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	task.CreatedAt = &createdAt
	task.UpdatedAt = &updatedAt

	return task, nil
}

func (r *PostgresTaskRepository) GetAll(subject string) ([]*model.Task, error) {
	query := `
        SELECT id, title, description, due_date, done, created_at, updated_at
        FROM tasks
        WHERE subject = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Done,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		createdAtStr := createdAt.Format(time.RFC3339)
		updatedAtStr := updatedAt.Format(time.RFC3339)
		task.CreatedAt = &createdAtStr
		task.UpdatedAt = &updatedAtStr
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(id string, subject string) (*model.Task, error) {
	query := `
        SELECT id, title, description, due_date, done, created_at, updated_at
        FROM tasks
        WHERE id = $1 AND subject = $2
    `

	task := &model.Task{}
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(query, id, subject).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Done,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	createdAtStr := createdAt.Format(time.RFC3339)
	updatedAtStr := updatedAt.Format(time.RFC3339)
	task.CreatedAt = &createdAtStr
	task.UpdatedAt = &updatedAtStr

	return task, nil
}

func (r *PostgresTaskRepository) Update(id string, input model.UpdateTaskInput, subject string) (*model.Task, error) {
	task, err := r.GetByID(id, subject)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = input.Description
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.Done != nil {
		task.Done = *input.Done
	}
	now := time.Now()
	updatedAt := now.Format(time.RFC3339)

	query := `
        UPDATE tasks
        SET title = $1, description = $2, due_date = $3, done = $4, updated_at = $5
        WHERE id = $6 AND subject = $7
        RETURNING id, title, description, due_date, done, created_at, updated_at
    `

	var createdAt time.Time
	err = r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Done,
		now,
		id,
		subject,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Done,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	createdAtStr := createdAt.Format(time.RFC3339)
	task.CreatedAt = &createdAtStr
	task.UpdatedAt = &updatedAt

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
