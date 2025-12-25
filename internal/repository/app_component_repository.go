package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vhvcorp/go-system-config-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AppComponentRepository handles app component data access
type AppComponentRepository struct {
	collection *mongo.Collection
}

// NewAppComponentRepository creates a new app component repository
func NewAppComponentRepository(db *mongo.Database) *AppComponentRepository {
	collection := db.Collection("app_components")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tenant_id", Value: 1},
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "status", Value: 1}},
		},
	}
	
	_, _ = collection.Indexes().CreateMany(ctx, indexes)
	
	return &AppComponentRepository{collection: collection}
}

// Create creates a new app component
func (r *AppComponentRepository) Create(ctx context.Context, component *domain.AppComponent) error {
	component.CreatedAt = time.Now()
	component.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to create app component: %w", err)
	}
	
	component.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds an app component by ID
func (r *AppComponentRepository) FindByID(ctx context.Context, id string) (*domain.AppComponent, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid app component ID: %w", err)
	}
	
	var component domain.AppComponent
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&component)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find app component: %w", err)
	}
	return &component, nil
}

// FindByCode finds an app component by code and tenant
func (r *AppComponentRepository) FindByCode(ctx context.Context, tenantID, code string) (*domain.AppComponent, error) {
	var component domain.AppComponent
	err := r.collection.FindOne(ctx, bson.M{
		"tenant_id": tenantID,
		"code":      code,
	}).Decode(&component)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find app component: %w", err)
	}
	return &component, nil
}

// List lists app components with pagination
func (r *AppComponentRepository) List(ctx context.Context, tenantID string, page, perPage int) ([]*domain.AppComponent, int64, error) {
	filter := bson.M{"tenant_id": tenantID}
	
	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count app components: %w", err)
	}
	
	// Find with pagination
	opts := options.Find().
		SetSkip(int64((page - 1) * perPage)).
		SetLimit(int64(perPage)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list app components: %w", err)
	}
	defer cursor.Close(ctx)
	
	var components []*domain.AppComponent
	if err = cursor.All(ctx, &components); err != nil {
		return nil, 0, fmt.Errorf("failed to decode app components: %w", err)
	}
	
	return components, total, nil
}

// Update updates an app component
func (r *AppComponentRepository) Update(ctx context.Context, component *domain.AppComponent) error {
	component.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": component.ID},
		bson.M{"$set": component},
	)
	if err != nil {
		return fmt.Errorf("failed to update app component: %w", err)
	}
	return nil
}

// Delete deletes an app component
func (r *AppComponentRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid app component ID: %w", err)
	}
	
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete app component: %w", err)
	}
	return nil
}
