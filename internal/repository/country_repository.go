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

// CountryRepository handles country data access
type CountryRepository struct {
	collection *mongo.Collection
}

// NewCountryRepository creates a new country repository
func NewCountryRepository(db *mongo.Database) *CountryRepository {
	collection := db.Collection("countries")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "code", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
	}
	
	_, _ = collection.Indexes().CreateMany(ctx, indexes)
	
	return &CountryRepository{collection: collection}
}

// Create creates a new country
func (r *CountryRepository) Create(ctx context.Context, country *domain.Country) error {
	country.CreatedAt = time.Now()
	country.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, country)
	if err != nil {
		return fmt.Errorf("failed to create country: %w", err)
	}
	
	country.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByCode finds a country by code
func (r *CountryRepository) FindByCode(ctx context.Context, code string) (*domain.Country, error) {
	var country domain.Country
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&country)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find country: %w", err)
	}
	return &country, nil
}

// List lists all countries
func (r *CountryRepository) List(ctx context.Context, page, perPage int) ([]*domain.Country, int64, error) {
	filter := bson.M{"status": "active"}
	
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count countries: %w", err)
	}
	
	opts := options.Find().
		SetSkip(int64((page - 1) * perPage)).
		SetLimit(int64(perPage)).
		SetSort(bson.D{{Key: "code", Value: 1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list countries: %w", err)
	}
	defer cursor.Close(ctx)
	
	var countries []*domain.Country
	if err = cursor.All(ctx, &countries); err != nil {
		return nil, 0, fmt.Errorf("failed to decode countries: %w", err)
	}
	
	return countries, total, nil
}

// Update updates a country
func (r *CountryRepository) Update(ctx context.Context, country *domain.Country) error {
	country.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": country.ID},
		bson.M{"$set": country},
	)
	if err != nil {
		return fmt.Errorf("failed to update country: %w", err)
	}
	return nil
}

// Delete deletes a country
func (r *CountryRepository) Delete(ctx context.Context, code string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"code": code})
	if err != nil {
		return fmt.Errorf("failed to delete country: %w", err)
	}
	return nil
}
