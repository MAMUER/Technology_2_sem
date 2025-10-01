package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "example.com/pz3-http/internal/storage"
)

func TestCreateTask(t *testing.T) {
    store := storage.NewMemoryStore()
    h := NewHandlers(store)

    payload := `{"title":"Test task"}`
    req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    h.CreateTask(w, req)

    res := w.Result()
    if res.StatusCode != http.StatusCreated {
        t.Fatalf("expected status 201, got %d", res.StatusCode)
    }

    var tsk storage.Task
    if err := json.NewDecoder(res.Body).Decode(&tsk); err != nil {
        t.Fatal("failed to decode response:", err)
    }
    if tsk.Title != "Test task" || tsk.Done {
        t.Fatalf("unexpected task data: %+v", tsk)
    }
}
