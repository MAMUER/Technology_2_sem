package sanitize

import (
	"html"
	"strings"
)

// SanitizeText удаляет потенциально опасные символы из текста
func SanitizeText(input string) string {
	if input == "" {
		return input
	}

	// Экранируем HTML специальные символы
	sanitized := html.EscapeString(input)

	// Удаляем потенциально опасные последовательности
	sanitized = strings.ReplaceAll(sanitized, "javascript:", "")
	sanitized = strings.ReplaceAll(sanitized, "data:", "")
	sanitized = strings.ReplaceAll(sanitized, "vbscript:", "")

	return sanitized
}

// SanitizeHTML очищает HTML (более агрессивно)
func SanitizeHTML(input string) string {
	if input == "" {
		return input
	}

	// Удаляем все HTML теги
	sanitized := stripHTMLTags(input)

	// Экранируем оставшиеся специальные символы
	sanitized = html.EscapeString(sanitized)

	return sanitized
}

// stripHTMLTags удаляет HTML теги из строки
func stripHTMLTags(input string) string {
	var result strings.Builder
	inTag := false

	for _, r := range input {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ValidateAndSanitizeDescription проверяет и очищает описание задачи
func ValidateAndSanitizeDescription(description string) (string, error) {
	if description == "" {
		return description, nil
	}

	// Проверяем длину
	if len(description) > 1000 {
		return "", ErrDescriptionTooLong
	}

	// Очищаем
	sanitized := SanitizeHTML(description)

	return sanitized, nil
}

// Определяем ошибки
var (
	ErrDescriptionTooLong = &ValidationError{"description too long (max 1000 characters)"}
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
