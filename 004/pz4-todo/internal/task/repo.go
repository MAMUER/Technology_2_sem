package task

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]*Task
	filename string
}

func NewRepo(filename string) *Repo {
	repo := &Repo{
        items:    make(map[int64]*Task),
        filename: filename,
    }
    repo.LoadFromFile(filename)
	return repo
	//return &Repo{items: make(map[int64]*Task)}
}

func (r *Repo) List() []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		out = append(out, t)
	}
	return out
}

func (r *Repo) Get(id int64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *Repo) Create(title string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	now := time.Now()
	t := &Task{ID: r.seq, Title: title, CreatedAt: now, UpdatedAt: now, Done: false}
	r.items[t.ID] = t
	_ = r.SaveToFile(r.filename)
	return t
}

func (r *Repo) Update(id int64, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()
	_ = r.SaveToFile(r.filename)
	return t, nil
}

func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	_ = r.SaveToFile(r.filename)
	return nil
}

// Новый метод загрузки из файла
func (r *Repo) LoadFromFile(filename string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // нет файла — просто пустой репозиторий
		}
		return err
	}
	defer file.Close()

	var tasks []*Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		return err
	}

	r.items = make(map[int64]*Task, len(tasks))
	r.seq = 0
	for _, t := range tasks {
		r.items[t.ID] = t
		if t.ID > r.seq {
			r.seq = t.ID
		}
	}
	return nil
}

// Новый метод сохранения в файл
func (r *Repo) SaveToFile(filename string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		tasks = append(tasks, t)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
