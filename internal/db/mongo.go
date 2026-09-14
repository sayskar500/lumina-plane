package db

import (
	"context"
	"time"

	"github.com/sayskar500/lumina-plane/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewStore(client *mongo.Client) *Store {
	return &Store{
		client: client,
		db:     client.Database("lumina_plane"),
	}
}

func (s *Store) LogTokens(ctx context.Context, log models.TokenLog) error {
	collection := s.db.Collection("token_logs")
	_, err := collection.InsertOne(ctx, log)
	return err
}

// SavePrompt stores a new prompt version
func (s *Store) SavePrompt(ctx context.Context, p models.Prompt) error {
	collection := s.db.Collection("prompts")

	// Deactivate old prompts for this project
	filter := bson.M{"project_id": p.ProjectID, "is_active": true}
	update := bson.M{"$set": bson.M{"is_active": false}}
	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

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
