package db

import (
	"database/sql"
	"log"
)

func MustApplyMigrations(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		
		-- Функция для автоматического обновления updated_at
		CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
		BEGIN 
			NEW.updated_at = NOW(); 
			RETURN NEW; 
		END; 
		$$ LANGUAGE plpgsql;
		
		-- Триггер для автоматического обновления updated_at
		DROP TRIGGER IF EXISTS trg_notes_updated ON notes;
		CREATE TRIGGER trg_notes_updated BEFORE UPDATE ON notes
		FOR EACH ROW EXECUTE FUNCTION set_updated_at();
		
		-- Индексы для производительности
		CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_notes_title ON notes(title);
	`)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Database migrations applied successfully")
}

// CleanTestData очищает тестовые данные
func CleanTestData(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM notes")
	return err
}
