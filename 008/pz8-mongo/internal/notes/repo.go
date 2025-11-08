package notes

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("note not found")

type Repo struct {
	col *mongo.Collection
}

func NewRepo(db *mongo.Database) (*Repo, error) {
	col := db.Collection("notes")
	textIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "content", Value: "text"},
		},
	}
	uniqueIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		textIndex,
		uniqueIndex,
	})
	if err != nil {
		return nil, err
	}

	return &Repo{col: col}, nil
}

// Create создает новую заметку
func (r *Repo) Create(ctx context.Context, title, content string) (Note, error) {
	now := time.Now()
	n := Note{Title: title, Content: content, CreatedAt: now, UpdatedAt: now}
	res, err := r.col.InsertOne(ctx, n)
	if err != nil {
		return Note{}, err
	}
	n.ID = res.InsertedID.(primitive.ObjectID)
	return n, nil
}

// ByID возвращает заметку по ID
func (r *Repo) ByID(ctx context.Context, idHex string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}
	var n Note
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&n); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return n, nil
}

// List возвращает список заметок с пагинацией и поиском
func (r *Repo) List(ctx context.Context, q string, limit, skip int64) ([]Note, error) {
	filter := bson.M{}
	if q != "" {
		filter["title"] = bson.M{"$regex": q, "$options": "i"}
	}
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, cur.Err()
}

// Update обновляет заметку
func (r *Repo) Update(ctx context.Context, idHex string, title, content *string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}

	set := bson.M{"updatedAt": time.Now()}
	if title != nil {
		set["title"] = *title
	}
	if content != nil {
		set["content"] = *content
	}

	after := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Note
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, after).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return updated, nil
}

// Delete удаляет заметку
func (r *Repo) Delete(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return ErrNotFound
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// TextSearch выполняет полнотекстовый поиск
func (r *Repo) TextSearch(ctx context.Context, query string, limit, skip int64) ([]Note, error) {
	filter := bson.M{}
	if query != "" {
		filter["$text"] = bson.M{"$search": query}
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var notes []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, cur.Err()
}

// TextSearchWithScore выполняет полнотекстовый поиск с релевантностью
func (r *Repo) TextSearchWithScore(ctx context.Context, query string, limit, skip int64) ([]NoteWithScore, error) {
	if query == "" {
		notes, err := r.List(ctx, "", limit, skip)
		if err != nil {
			return nil, err
		}

		var result []NoteWithScore
		for _, note := range notes {
			result = append(result, NoteWithScore{
				Note:  note,
				Score: 1.0,
			})
		}
		return result, nil
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"$text": bson.M{"$search": query}}}},
		{{Key: "$addFields", Value: bson.M{"score": bson.M{"$meta": "textScore"}}}},
		{{Key: "$sort", Value: bson.M{"score": bson.M{"$meta": "textScore"}, "createdAt": -1}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []NoteWithScore
	for cur.Next(ctx) {
		var result NoteWithScore
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, cur.Err()
}

// GetStats возвращает статистику по заметкам
func (r *Repo) GetStats(ctx context.Context) (StatsResponse, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.M{
			"contentLength": bson.M{"$strLenCP": "$content"},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":              nil,
			"totalNotes":       bson.M{"$sum": 1},
			"avgContentLength": bson.M{"$avg": "$contentLength"},
			"maxContentLength": bson.M{"$max": "$contentLength"},
			"minContentLength": bson.M{"$min": "$contentLength"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":              0,
			"totalNotes":       1,
			"avgContentLength": bson.M{"$round": bson.A{"$avgContentLength", 2}},
			"maxContentLength": 1,
			"minContentLength": 1,
		}}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return StatsResponse{}, err
	}
	defer cur.Close(ctx)

	var results []StatsResponse
	for cur.Next(ctx) {
		var stats StatsResponse
		if err := cur.Decode(&stats); err != nil {
			return StatsResponse{}, err
		}
		results = append(results, stats)
	}

	if len(results) == 0 {
		return StatsResponse{
			TotalNotes:       0,
			AvgContentLength: 0,
			MaxContentLength: 0,
			MinContentLength: 0,
		}, nil
	}

	return results[0], cur.Err()
}

// GetStatsByDay возвращает статистику по дням (количество заметок в день)
func (r *Repo) GetStatsByDay(ctx context.Context, days int) ([]DayStats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"createdAt": bson.M{
				"$gte": time.Now().AddDate(0, 0, -days),
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$createdAt"},
				"month": bson.M{"$month": "$createdAt"},
				"day":   bson.M{"$dayOfMonth": "$createdAt"},
			},
			"count": bson.M{"$sum": 1},
			"date":  bson.M{"$first": "$createdAt"},
		}}},
		{{Key: "$sort", Value: bson.M{"date": 1}}},
		{{Key: "$project", Value: bson.M{
			"_id":   0,
			"date":  bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$date"}},
			"count": 1,
		}}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []DayStats
	for cur.Next(ctx) {
		var stats DayStats
		if err := cur.Decode(&stats); err != nil {
			return nil, err
		}
		results = append(results, stats)
	}

	return results, cur.Err()
}
