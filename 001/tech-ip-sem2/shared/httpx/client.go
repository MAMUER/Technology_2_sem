package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client обёртка над HTTP-клиентом с таймаутами
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создаёт нового клиента
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do выполняет HTTP-запрос с прокидыванием контекста и RequestID
func (c *Client) Do(ctx context.Context, method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Устанавливаем заголовки
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Прокидываем существующие заголовки
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Прокидываем RequestID из контекста, если есть
	if requestID, ok := ctx.Value("requestID").(string); ok && requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	return c.httpClient.Do(req)
}
