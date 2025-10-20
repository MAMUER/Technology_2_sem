package notes

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title"         json:"title"`
	Content   string             `bson:"content"       json:"content"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"     json:"updatedAt"`
}

// NoteWithScore для полнотекстового поиска с релевантностью
type NoteWithScore struct {
	Note  `bson:",inline"`
	Score float64 `bson:"score" json:"score"`
}

// StatsResponse для агрегационной статистики
type StatsResponse struct {
	TotalNotes       int64   `bson:"totalNotes"   json:"totalNotes"`
	AvgContentLength float64 `bson:"avgContentLength" json:"avgContentLength"`
	MaxContentLength int     `bson:"maxContentLength" json:"maxContentLength"`
	MinContentLength int     `bson:"minContentLength" json:"minContentLength"`
}

// DayStats для статистики по дням
type DayStats struct {
	Date  string `bson:"date"  json:"date"`
	Count int    `bson:"count" json:"count"`
}

// SearchResponse для универсального ответа поиска
type SearchResponse struct {
	Notes   []NoteWithScore `json:"notes"`
	Total   int64           `json:"total"`
	Query   string          `json:"query,omitempty"`
	HasMore bool            `json:"hasMore"`
}
