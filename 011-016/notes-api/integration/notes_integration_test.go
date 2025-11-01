package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/notes-api/internal/core/service"
	httpx "example.com/notes-api/internal/httpapi"
	"example.com/notes-api/internal/httpapi/handlers"
	"example.com/notes-api/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// testContainer создает и запускает контейнер PostgreSQL для тестов
func startTestDB(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, error) {
	// Запускаем PostgreSQL контейнер
	container, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("notes_test"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Получаем DSN для подключения
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Подключаемся к БД
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Применяем миграции
	if err := applyMigrations(ctx, dbPool); err != nil {
		return nil, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return container, dbPool, nil
}

// applyMigrations применяет схему БД
func applyMigrations(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS notes (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		
		CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_notes_title ON notes(title);
	`)
	return err
}

// cleanTestData очищает данные между тестами
func cleanTestData(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, "TRUNCATE TABLE notes RESTART IDENTITY CASCADE")
	return err
}

// newTestServer создает тестовый HTTP сервер
func newTestServer(db *pgxpool.Pool) *httptest.Server {
	// Создаем сервисы
	noteRepo := repo.NewNoteRepoPostgres(db)
	noteService := service.NewNoteService(noteRepo)

	// Создаем хендлер и роутер
	h := &handlers.Handler{Service: noteService}
	router := httpx.NewRouter(h)

	return httptest.NewServer(router)
}

func TestCreateNote_Success(t *testing.T) {
    ctx := context.Background()

    container, dbPool, err := startTestDB(ctx)
    if err != nil {
        t.Fatalf("Failed to start test DB: %v", err)
    }
    defer container.Terminate(ctx)
    defer dbPool.Close()

    if err := cleanTestData(ctx, dbPool); err != nil {
        t.Fatalf("Failed to clean test data: %v", err)
    }

    server := newTestServer(dbPool)
    defer server.Close()

    // Тело запроса
    noteData := map[string]string{
        "title":   "Integration Test Note", 
        "content": "This is a test note created by integration test",
    }

    jsonData, err := json.Marshal(noteData)
    if err != nil {
        t.Fatalf("Failed to marshal note data: %v", err)
    }

    // Выполняем запрос
    resp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        t.Fatalf("Failed to create note: %v", err)
    }
    defer resp.Body.Close()

    // Проверяем статус
    if resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        t.Fatalf("Expected status 201, got %d. Response: %s", resp.StatusCode, string(body))
    }

    // 🔥 ИСПРАВЛЕНИЕ: Ожидаем только ID в ответе
    var response struct {
        ID int `json:"id"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    // Проверяем что ID создан
    if response.ID == 0 {
        t.Error("Created note should have a non-zero ID")
    }

    // 🔥 ДОПОЛНИТЕЛЬНАЯ ПРОВЕРКА: Получаем созданную заметку чтобы убедиться в данных
    getResp, err := http.Get(fmt.Sprintf("%s/api/v1/notes/%d", server.URL, response.ID))
    if err != nil {
        t.Fatalf("Failed to get created note: %v", err)
    }
    defer getResp.Body.Close()

    if getResp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status 200 for get, got %d", getResp.StatusCode)
    }

    var retrievedNote map[string]interface{}
    if err := json.NewDecoder(getResp.Body).Decode(&retrievedNote); err != nil {
        t.Fatalf("Failed to decode retrieved note: %v", err)
    }

    // Теперь проверяем данные
    if retrievedNote["title"] != noteData["title"] {
        t.Errorf("Expected title %s, got %s", noteData["title"], retrievedNote["title"])
    }
    if retrievedNote["content"] != noteData["content"] {
        t.Errorf("Expected content %s, got %s", noteData["content"], retrievedNote["content"])
    }
}

func TestGetNote_Success(t *testing.T) {
	ctx := context.Background()

	container, dbPool, err := startTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to start test DB: %v", err)
	}
	defer container.Terminate(ctx)
	defer dbPool.Close()

	if err := cleanTestData(ctx, dbPool); err != nil {
		t.Fatalf("Failed to clean test data: %v", err)
	}

	server := newTestServer(dbPool)
	defer server.Close()

	// Сначала создаем заметку
	createData := map[string]string{
		"title":   "Note to Retrieve",
		"content": "Content for retrieval test",
	}

	createJSON, _ := json.Marshal(createData)
	createResp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(createJSON))
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

	// Проверяем статус
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Проверяем данные
	var retrievedNote map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&retrievedNote); err != nil {
		t.Fatalf("Failed to decode note: %v", err)
	}

	if retrievedNote["title"] != createData["title"] {
		t.Errorf("Expected title %s, got %s", createData["title"], retrievedNote["title"])
	}
	if retrievedNote["content"] != createData["content"] {
		t.Errorf("Expected content %s, got %s", createData["content"], retrievedNote["content"])
	}
}

func TestGetNote_NotFound(t *testing.T) {
	ctx := context.Background()

	container, dbPool, err := startTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to start test DB: %v", err)
	}
	defer container.Terminate(ctx)
	defer dbPool.Close()

	if err := cleanTestData(ctx, dbPool); err != nil {
		t.Fatalf("Failed to clean test data: %v", err)
	}

	server := newTestServer(dbPool)
	defer server.Close()

	// Пытаемся получить несуществующую заметку
	resp, err := http.Get(server.URL + "/api/v1/notes/9999")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем что получили 404
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status 404 for non-existent note, got %d", resp.StatusCode)
	}
}

func TestUpdateNote_Success(t *testing.T) {
	ctx := context.Background()

	container, dbPool, err := startTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to start test DB: %v", err)
	}
	defer container.Terminate(ctx)
	defer dbPool.Close()

	if err := cleanTestData(ctx, dbPool); err != nil {
		t.Fatalf("Failed to clean test data: %v", err)
	}

	server := newTestServer(dbPool)
	defer server.Close()

	// Создаем заметку для обновления
	createData := map[string]string{
		"title":   "Original Title",
		"content": "Original Content",
	}

	createJSON, _ := json.Marshal(createData)
	createResp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(createJSON))
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

	// Проверяем успешное обновление
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
	}
}

func TestDeleteNote_Success(t *testing.T) {
	ctx := context.Background()

	container, dbPool, err := startTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to start test DB: %v", err)
	}
	defer container.Terminate(ctx)
	defer dbPool.Close()

	if err := cleanTestData(ctx, dbPool); err != nil {
		t.Fatalf("Failed to clean test data: %v", err)
	}

	server := newTestServer(dbPool)
	defer server.Close()

	// Создаем заметку для удаления
	createData := map[string]string{
		"title":   "Note to Delete",
		"content": "This note will be deleted",
	}

	createJSON, _ := json.Marshal(createData)
	createResp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(createJSON))
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

	// Проверяем успешное удаление
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status 200 or 204, got %d", resp.StatusCode)
	}

	// Проверяем что заметка действительно удалена
	getResp, err := http.Get(fmt.Sprintf("%s/api/v1/notes/%d", server.URL, noteID))
	if err != nil {
		t.Fatalf("Failed to check deleted note: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status 404 after deletion, got %d", getResp.StatusCode)
	}
}

func TestGetAllNotes_Success(t *testing.T) {
	ctx := context.Background()

	container, dbPool, err := startTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to start test DB: %v", err)
	}
	defer container.Terminate(ctx)
	defer dbPool.Close()

	if err := cleanTestData(ctx, dbPool); err != nil {
		t.Fatalf("Failed to clean test data: %v", err)
	}

	server := newTestServer(dbPool)
	defer server.Close()

	// Создаем несколько заметок
	notes := []map[string]string{
		{"title": "Note 1", "content": "Content 1"},
		{"title": "Note 2", "content": "Content 2"},
	}

	for _, noteData := range notes {
		jsonData, _ := json.Marshal(noteData)
		resp, err := http.Post(server.URL+"/api/v1/notes", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create note: %v", err)
		}
		resp.Body.Close()
	}

	// Получаем все заметки
	resp, err := http.Get(server.URL + "/api/v1/notes")
	if err != nil {
		t.Fatalf("Failed to get notes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var allNotes []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&allNotes); err != nil {
		t.Fatalf("Failed to decode notes: %v", err)
	}

	// Проверяем что получили все созданные заметки
	if len(allNotes) != len(notes) {
		t.Errorf("Expected %d notes, got %d", len(notes), len(allNotes))
	}
}
