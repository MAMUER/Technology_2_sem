package service

import (
	"example.com/notes-api/internal/core"
	"time"
)

// Определяем интерфейс репозитория в этом пакете
type NoteRepository interface {
	Create(n core.Note) (int64, error)
	GetAll() ([]*core.Note, error)
	GetByID(id int64) (*core.Note, error)
	Update(updatedNote *core.Note) (*core.Note, error)
	Delete(id int64) error
}

type NoteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(title, content string) (int64, error) {
	note := core.Note{
		Title:   title,
		Content: content,
	}
	return s.repo.Create(note)
}

func (s *NoteService) GetAllNotes() ([]*core.Note, error) {
	return s.repo.GetAll()
}

func (s *NoteService) GetNoteByID(id int64) (*core.Note, error) {
	return s.repo.GetByID(id)
}

func (s *NoteService) UpdateNote(id int64, title, content string) (*core.Note, error) {
	note := &core.Note{
		ID:      id,
		Title:   title,
		Content: content,
	}
	return s.repo.Update(note)
}

func (s *NoteService) DeleteNote(id int64) error {
	return s.repo.Delete(id)
}

// GetAllNotesWithPagination - keyset пагинация
func (s *NoteService) GetAllNotesWithPagination(cursorTime time.Time, cursorID int64, limit int) ([]*core.Note, error) {
	// Приводим к конкретному типу для доступа к методам пагинации
	if repo, ok := s.repo.(interface {
		GetAllWithPagination(cursorTime time.Time, cursorID int64, limit int) ([]*core.Note, error)
	}); ok {
		return repo.GetAllWithPagination(cursorTime, cursorID, limit)
	}

	// Fallback на обычный GetAll если репозиторий не поддерживает пагинацию
	return s.repo.GetAll()
}

// GetNotesByIDs - батчинг для получения нескольких заметок
func (s *NoteService) GetNotesByIDs(ids []int64) ([]*core.Note, error) {
	if repo, ok := s.repo.(interface {
		GetByIDs(ids []int64) ([]*core.Note, error)
	}); ok {
		return repo.GetByIDs(ids)
	}

	// Fallback - последовательные запросы (неэффективно)
	var notes []*core.Note
	for _, id := range ids {
		note, err := s.repo.GetByID(id)
		if err != nil {
			continue // или вернуть ошибку
		}
		notes = append(notes, note)
	}
	return notes, nil
}
