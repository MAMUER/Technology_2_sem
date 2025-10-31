package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/notes-api/internal/core/service"
	httpx "example.com/notes-api/internal/httpapi"
	"example.com/notes-api/internal/httpapi/handlers"
	"example.com/notes-api/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// newTestServer создает тестовый сервер с чистой БД
func newTestServer(t *testing.T) (*httptest.Server, *pgxpool.Pool) {
	t.Helper()

	// Используем тестовую БД
	dsn := getTestDSN()
	
	// Создаем pgxpool вместо sql.DB
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Применяем миграции
	if err := applyMigrations(dbPool); err != nil {
		t.Fatalf("Failed to apply migrations: %v", err)
	}

	// Очищаем данные перед тестом
	if err := cleanTestData(dbPool); err != nil {
		t.Fatalf("Failed to clean test data: %v", err)
	}

	// Создаем сервисы
	noteRepo := repo.NewNoteRepoPostgres(dbPool)
	noteService := service.NewNoteService(noteRepo)

	// Создаем хендлер
	h := &handlers.Handler{Service: noteService}

	// Создаем роутер
	router := httpx.NewRouter(h)

	// Создаем тестовый сервер
	server := httptest.NewServer(router)

	return server, dbPool
}

// applyMigrations применяет миграции для pgxpool
func applyMigrations(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS notes (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		
		-- Индексы для производительности
		CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_notes_title ON notes(title);
	`)
	return err
}

// cleanTestData очищает тестовые данные
func cleanTestData(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), "DELETE FROM notes")
	return err
}

// getTestDSN возвращает DSN для тестовой БД
func getTestDSN() string {
	if dsn := os.Getenv("TEST_DATABASE_URL"); dsn != "" {
		return dsn
	}
	// Fallback для локальной разработки
	return "postgres://teacher_app:secure_password_123@localhost:5434/notes?sslmode=disable"
}

func TestCreateAndGetNote(t *testing.T) {
	server, dbPool := newTestServer(t)
	defer server.Close()
	defer dbPool.Close()

	// 1. Тест создания заметки
	t.Run("Create Note", func(t *testing.T) {
		noteData := map[string]string{
			"title":   "Integration Test Note",
			"content": "This is a test note created by integration test",
		}

		jsonData, err := json.Marshal(noteData)
		if err != nil {
			t.Fatalf("Failed to marshal note data: %v", err)
		}

		resp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d. Response: %s", resp.StatusCode, string(body))
		}

		var createdNote map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Проверяем что ID был присвоен
		if _, exists := createdNote["id"]; !exists {
			t.Error("Created note should have an ID")
		}

		if createdNote["title"] != noteData["title"] {
			t.Errorf("Expected title %s, got %s", noteData["title"], createdNote["title"])
		}
	})

	// 2. Тест получения списка заметок
	t.Run("Get Notes List", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/notes")
		if err != nil {
			t.Fatalf("Failed to get notes list: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var notes []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
			t.Fatalf("Failed to decode notes list: %v", err)
		}

		// Проверяем что есть хотя бы одна заметка (созданная в предыдущем тесте)
		if len(notes) == 0 {
			t.Error("Expected at least one note in the list")
		}
	})

	// 3. Тест получения конкретной заметки
	t.Run("Get Specific Note", func(t *testing.T) {
		// Сначала создаем заметку чтобы получить её ID
		noteData := map[string]string{
			"title":   "Note to Retrieve",
			"content": "Content for retrieval test",
		}

		jsonData, _ := json.Marshal(noteData)
		createResp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}

		var createdNote map[string]interface{}
		json.NewDecoder(createResp.Body).Decode(&createdNote)
		createResp.Body.Close()

		noteID := int(createdNote["id"].(float64))

		// Теперь получаем созданную заметку
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/notes/%d", server.URL, noteID))
		if err != nil {
			t.Fatalf("Failed to get note: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var retrievedNote map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&retrievedNote); err != nil {
			t.Fatalf("Failed to decode note: %v", err)
		}

		if retrievedNote["title"] != noteData["title"] {
			t.Errorf("Expected title %s, got %s", noteData["title"], retrievedNote["title"])
		}
	})
}

func TestUpdateNote(t *testing.T) {
	server, dbPool := newTestServer(t)
	defer server.Close()
	defer dbPool.Close()

	// Создаем заметку для обновления
	noteData := map[string]string{
		"title":   "Original Title",
		"content": "Original Content",
	}

	jsonData, _ := json.Marshal(noteData)
	createResp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	var createdNote map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createdNote)
	createResp.Body.Close()

	noteID := int(createdNote["id"].(float64))

	// Обновляем заметку
	updateData := map[string]string{
		"title": "Updated Title",
	}

	updateJSON, _ := json.Marshal(updateData)
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/notes/%d", server.URL, noteID), bytes.NewBuffer(updateJSON))
	if err != nil {
		t.Fatalf("Failed to create update request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to update note: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
	}
}

func TestDeleteNote(t *testing.T) {
	server, dbPool := newTestServer(t)
	defer server.Close()
	defer dbPool.Close()

	// Создаем заметку для удаления
	noteData := map[string]string{
		"title":   "Note to Delete",
		"content": "This note will be deleted",
	}

	jsonData, _ := json.Marshal(noteData)
	createResp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	var createdNote map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createdNote)
	createResp.Body.Close()

	noteID := int(createdNote["id"].(float64))

	// Удаляем заметку
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/notes/%d", server.URL, noteID), nil)
	if err != nil {
		t.Fatalf("Failed to create delete request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete note: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status 200 or 204, got %d", resp.StatusCode)
	}
}

// Тест для проверки пагинации
func TestNotesPagination(t *testing.T) {
	server, dbPool := newTestServer(t)
	defer server.Close()
	defer dbPool.Close()

	// Создаем несколько заметок для тестирования пагинации
	notes := []map[string]string{
		{"title": "Note 1", "content": "Content 1"},
		{"title": "Note 2", "content": "Content 2"},
		{"title": "Note 3", "content": "Content 3"},
	}

	for _, noteData := range notes {
		jsonData, _ := json.Marshal(noteData)
		resp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}
		resp.Body.Close()
	}

	// Тестируем пагинацию
	resp, err := http.Get(server.URL + "/api/v1/notes/paginated?limit=2")
	if err != nil {
		t.Fatalf("Failed to get paginated notes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var paginatedNotes []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&paginatedNotes); err != nil {
		t.Fatalf("Failed to decode paginated notes: %v", err)
	}

	if len(paginatedNotes) > 2 {
		t.Errorf("Expected max 2 notes with pagination, got %d", len(paginatedNotes))
	}
}