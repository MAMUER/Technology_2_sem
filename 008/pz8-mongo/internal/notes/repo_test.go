package notes

import (
	"context"
	"os"
	"testing"
	"time"

	"example.com/pz8-mongo/internal/db"
)

// TestCreateAndGet тестирует создание и получение заметки
func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()

	// Используем тестовую базу данных с уникальным именем
	dbName := "pz8_test_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Очистка после теста
	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Тест создания
	created, err := repo.Create(ctx, "Test Title", "Test Content")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	// Проверяем что ID установлен
	if created.ID.IsZero() {
		t.Error("Expected note to have ID, got zero value")
	}

	// Тест получения по ID
	got, err := repo.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get note by ID: %v", err)
	}

	// Проверяем поля
	if got.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", got.Title)
	}
	if got.Content != "Test Content" {
		t.Errorf("Expected content 'Test Content', got '%s'", got.Content)
	}
	if got.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if got.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

// TestCreateDuplicateTitle тестирует создание заметки с дублирующимся заголовком
func TestCreateDuplicateTitle(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_duplicate_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Первая заметка - должна создаться успешно
	_, err = repo.Create(ctx, "Duplicate Title", "Content 1")
	if err != nil {
		t.Fatalf("Failed to create first note: %v", err)
	}

	// Вторая заметка с тем же заголовком - должна вернуть ошибку
	_, err = repo.Create(ctx, "Duplicate Title", "Content 2")
	if err == nil {
		t.Error("Expected error for duplicate title, got nil")
	}
}

// TestByID_NotFound тестирует получение несуществующей заметки
func TestByID_NotFound(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_notfound_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Пытаемся получить несуществующую заметку
	_, err = repo.ByID(ctx, "507f1f77bcf86cd799439011") // valid ObjectID format
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestByID_InvalidID тестирует получение с невалидным ID
func TestByID_InvalidID(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_invalid_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Пытаемся получить с невалидным ID
	_, err = repo.ByID(ctx, "invalid-id")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound for invalid ID, got %v", err)
	}
}

// TestList тестирует получение списка заметок
func TestList(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_list_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Создаем несколько заметок
	notes := []struct {
		title   string
		content string
	}{
		{"Note 1", "Content 1"},
		{"Note 2", "Content 2"},
		{"Another Note", "Content 3"},
	}

	for _, n := range notes {
		_, err := repo.Create(ctx, n.title, n.content)
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}
	}

	// Тестируем получение всех заметок
	allNotes, err := repo.List(ctx, "", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notes: %v", err)
	}

	if len(allNotes) != 3 {
		t.Errorf("Expected 3 notes, got %d", len(allNotes))
	}

	// Тестируем поиск по заголовку
	searchNotes, err := repo.List(ctx, "Note", 10, 0)
	if err != nil {
		t.Fatalf("Failed to search notes: %v", err)
	}

	if len(searchNotes) != 3 {
		t.Errorf("Expected 3 notes with 'Note' in title, got %d", len(searchNotes))
	}

	// Тестируем пагинацию
	paginatedNotes, err := repo.List(ctx, "", 2, 0)
	if err != nil {
		t.Fatalf("Failed to list paginated notes: %v", err)
	}

	if len(paginatedNotes) != 2 {
		t.Errorf("Expected 2 notes with limit 2, got %d", len(paginatedNotes))
	}
}

// TestUpdate тестирует обновление заметки
func TestUpdate(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_update_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Создаем заметку
	created, err := repo.Create(ctx, "Original Title", "Original Content")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	// Обновляем только content
	newContent := "Updated Content"
	updated, err := repo.Update(ctx, created.ID.Hex(), nil, &newContent)
	if err != nil {
		t.Fatalf("Failed to update note: %v", err)
	}

	// Проверяем что title остался прежним, а content обновился
	if updated.Title != "Original Title" {
		t.Errorf("Expected title to remain 'Original Title', got '%s'", updated.Title)
	}
	if updated.Content != "Updated Content" {
		t.Errorf("Expected content 'Updated Content', got '%s'", updated.Content)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

// TestDelete тестирует удаление заметки
func TestDelete(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_delete_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Создаем заметку
	created, err := repo.Create(ctx, "To Delete", "Content")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	// Удаляем заметку
	err = repo.Delete(ctx, created.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to delete note: %v", err)
	}

	// Проверяем что заметка больше не существует
	_, err = repo.ByID(ctx, created.ID.Hex())
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}
}

// TestTextSearch тестирует полнотекстовый поиск
func TestTextSearch(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_textsearch_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Создаем тестовые заметки
	notes := []struct {
		title   string
		content string
	}{
		{"Go programming language", "Go is awesome for backend development"},
		{"JavaScript tutorial", "Learn JavaScript for web development"},
		{"Advanced Go patterns", "Concurrency patterns in Go"},
	}

	for _, n := range notes {
		_, err := repo.Create(ctx, n.title, n.content)
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}
	}

	// Даем время на создание индексов
	time.Sleep(2 * time.Second)

	// Тестируем поиск по слову "Go"
	results, err := repo.TextSearch(ctx, "Go", 10, 0)
	if err != nil {
		t.Fatalf("Failed to text search: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 notes with 'Go', got %d", len(results))
	}
}

// TestGetStats тестирует агрегационную статистику
func TestGetStats(t *testing.T) {
	ctx := context.Background()
	dbName := "pz8_test_stats_" + time.Now().Format("20060102150405")
	mongoURI := getTestMongoURI()

	deps, err := db.ConnectMongo(ctx, mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	t.Cleanup(func() {
		deps.Database.Drop(ctx)
		deps.Client.Disconnect(ctx)
	})

	repo, err := NewRepo(deps.Database)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Создаем заметки с разной длиной контента
	_, err = repo.Create(ctx, "Note 1", "Short")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	_, err = repo.Create(ctx, "Note 2", "Very long content here")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalNotes != 2 {
		t.Errorf("Expected 2 total notes, got %d", stats.TotalNotes)
	}

	if stats.AvgContentLength <= 0 {
		t.Errorf("Expected positive avg content length, got %f", stats.AvgContentLength)
	}
}

// getTestMongoURI возвращает URI для тестовой MongoDB
func getTestMongoURI() string {
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}
	return "mongodb://root:secret@localhost:27017/?authSource=admin"
}
