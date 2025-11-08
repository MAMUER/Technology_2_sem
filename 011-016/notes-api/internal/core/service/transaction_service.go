package service

import (
	"context"
	"fmt"

	"example.com/notes-api/internal/core"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	db *pgxpool.Pool
}

func NewTransactionService(db *pgxpool.Pool) *TransactionService {
	return &TransactionService{db: db}
}

// CreateNoteWithTags - транзакция с несколькими операциями
func (s *TransactionService) CreateNoteWithTags(note core.Note, tags []string) (int64, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var noteID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO notes (title, content, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id`,
		note.Title, note.Content, note.CreatedAt,
	).Scan(&noteID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert note: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return noteID, nil
}
