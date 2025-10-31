// internal/repo/note_mem.go
package repo

import (
	"errors"
	"example.com/notes-api/internal/core"
	"sync"
)

type NoteRepoMem struct {
	mu    sync.Mutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	n.ID = r.next
	r.notes[n.ID] = &n
	return n.ID, nil
}

func (r *NoteRepoMem) GetAll() ([]*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	notes := make([]*core.Note, 0, len(r.notes))
	for _, note := range r.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *NoteRepoMem) GetByID(id int64) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	note, exists := r.notes[id]
	if !exists {
		return nil, errors.New("note not found")
	}
	return note, nil
}

func (r *NoteRepoMem) Update(updatedNote *core.Note) (*core.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.notes[updatedNote.ID]
	if !exists {
		return nil, errors.New("note not found")
	}
	r.notes[updatedNote.ID] = updatedNote
	return updatedNote, nil
}

func (r *NoteRepoMem) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.notes[id]
	if !exists {
		return errors.New("note not found")
	}
	delete(r.notes, id)
	return nil
}
