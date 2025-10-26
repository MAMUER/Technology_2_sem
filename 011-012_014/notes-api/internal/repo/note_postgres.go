package repo

import (
	"context"
	"example.com/notes-api/internal/core"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepoPostgres struct {
	db *pgxpool.Pool
}

func NewNoteRepoPostgres(db *pgxpool.Pool) *NoteRepoPostgres {
	return &NoteRepoPostgres{db: db}
}

// Подготовленные запросы
const (
	createQuery = `
		INSERT INTO notes (title, content, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	getAllQuery = `
		SELECT id, title, content, created_at 
		FROM notes 
		ORDER BY created_at DESC, id DESC`

	getByIDQuery = `
		SELECT id, title, content, created_at 
		FROM notes 
		WHERE id = $1`

	updateQuery = `
		UPDATE notes 
		SET title = $1, content = $2 
		WHERE id = $3 
		RETURNING id, title, content, created_at`

	deleteQuery = `
		DELETE FROM notes WHERE id = $1`

	// Keyset пагинация
	getAllKeysetQuery = `
		SELECT id, title, content, created_at 
		FROM notes 
		WHERE (created_at, id) < ($1, $2)
		ORDER BY created_at DESC, id DESC 
		LIMIT $3`

	// Батчинг для получения нескольких записей
	getByIDsQuery = `
		SELECT id, title, content, created_at 
		FROM notes 
		WHERE id = ANY($1) 
		ORDER BY created_at DESC, id DESC`
)

func (r *NoteRepoPostgres) Create(n core.Note) (int64, error) {
	var id int64
	err := r.db.QueryRow(context.Background(), createQuery, n.Title, n.Content, time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create note: %w", err)
	}
	return id, nil
}

func (r *NoteRepoPostgres) GetAll() ([]*core.Note, error) {
	return r.GetAllWithLimit(10)
}

// GetAllWithPagination - keyset пагинация
func (r *NoteRepoPostgres) GetAllWithPagination(cursorTime time.Time, cursorID int64, limit int) ([]*core.Note, error) {
	rows, err := r.db.Query(context.Background(), getAllKeysetQuery, cursorTime, cursorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes with pagination: %w", err)
	}
	defer rows.Close()

	return r.scanNotes(rows)
}

// GetAllWithLimit - для OFFSET пагинации (менее эффективно)
func (r *NoteRepoPostgres) GetAllWithLimit(limit int) ([]*core.Note, error) {
	query := getAllQuery
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all notes: %w", err)
	}
	defer rows.Close()

	return r.scanNotes(rows)
}

func (r *NoteRepoPostgres) GetByID(id int64) (*core.Note, error) {
	var note core.Note
	err := r.db.QueryRow(context.Background(), getByIDQuery, id).Scan(
		&note.ID, &note.Title, &note.Content, &note.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get note by id: %w", err)
	}
	return &note, nil
}

// GetByIDs - батчинг для получения нескольких записей за один запрос
func (r *NoteRepoPostgres) GetByIDs(ids []int64) ([]*core.Note, error) {
	rows, err := r.db.Query(context.Background(), getByIDsQuery, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes by ids: %w", err)
	}
	defer rows.Close()

	return r.scanNotes(rows)
}

func (r *NoteRepoPostgres) Update(updatedNote *core.Note) (*core.Note, error) {
	var note core.Note
	err := r.db.QueryRow(context.Background(), updateQuery,
		updatedNote.Title, updatedNote.Content, updatedNote.ID).Scan(
		&note.ID, &note.Title, &note.Content, &note.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}
	return &note, nil
}

func (r *NoteRepoPostgres) Delete(id int64) error {
	result, err := r.db.Exec(context.Background(), deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

// Вспомогательная функция для сканирования результатов
func (r *NoteRepoPostgres) scanNotes(rows pgx.Rows) ([]*core.Note, error) {
	var notes []*core.Note
	for rows.Next() {
		var note core.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return notes, nil
}
