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

	// Экранирование HTML специальных символов
	sanitized := html.EscapeString(input)

	// Удаление потенциально опасных последовательностей
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

	// Удаление всех HTML тегов
	sanitized := stripHTMLTags(input)

	// Экранирование оставшихся специальных символов
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

	// Проверка длины
	if len(description) > 1000 {
		return "", ErrDescriptionTooLong
	}

	// Очистка
	sanitized := SanitizeHTML(description)

	return sanitized, nil
}

// Определение ошибок
var (
	ErrDescriptionTooLong = &ValidationError{"description too long (max 1000 characters)"}
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
