package db

import (
	"context"
	"time"

	"github.com/sayskar500/lumina-plane/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DefaultDatabaseName is used when no database name is configured.
const DefaultDatabaseName = "lumina_plane"

type Store struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewStore binds a Store to the given database name. An empty name falls back
// to DefaultDatabaseName.
func NewStore(client *mongo.Client, dbName string) *Store {
	if dbName == "" {
		dbName = DefaultDatabaseName
	}
	return &Store{
		client: client,
		db:     client.Database(dbName),
	}
}

func (s *Store) LogTokens(ctx context.Context, log models.TokenLog) error {
	collection := s.db.Collection("token_logs")
	_, err := collection.InsertOne(ctx, log)
	return err
}

// SavePrompt stores a new prompt version. The version number is assigned
// atomically from a per-project counter, so concurrent saves cannot collide
// and any client-supplied version is ignored.
func (s *Store) SavePrompt(ctx context.Context, p models.Prompt) error {
	collection := s.db.Collection("prompts")

	// Deactivate the currently active prompts for this project.
	filter := bson.M{"project_id": p.ProjectID, "is_active": true}
	update := bson.M{"$set": bson.M{"is_active": false}}
	if _, err := collection.UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	// Atomically reserve the next version number for this project.
	counters := s.db.Collection("counters")
	var result struct {
		Seq int `bson:"seq"`
	}
	err := counters.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "prompt_version:" + p.ProjectID},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		return err
	}

	p.Version = result.Seq
	p.CreatedAt = time.Now()
	p.IsActive = true
	_, err = collection.InsertOne(ctx, p)
	return err
}

// GetActivePrompt retrieves the current active prompt for a project
func (s *Store) GetActivePrompt(ctx context.Context, projectID string) (*models.Prompt, error) {
	collection := s.db.Collection("prompts")
	var p models.Prompt
	filter := bson.M{"project_id": projectID, "is_active": true}

	err := collection.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
