package repository

import (
	"context"
	"time"

	"github.com/vhvplatform/go-system-config-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	secretCollection          = "secrets"
	secretAccessLogCollection = "secret_access_log"
)

// SecretRepository handles secret data persistence
type SecretRepository struct {
	db *mongo.Database
}

// NewSecretRepository creates a new secret repository
func NewSecretRepository(db *mongo.Database) *SecretRepository {
	return &SecretRepository{db: db}
}

// Create creates a new secret
func (r *SecretRepository) Create(ctx context.Context, secret *domain.Secret) error {
	secret.CreatedAt = time.Now()
	secret.UpdatedAt = time.Now()
	secret.Version = 1
	secret.AccessCount = 0

	result, err := r.db.Collection(secretCollection).InsertOne(ctx, secret)
	if err != nil {
		return err
	}
	secret.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a secret by ID
func (r *SecretRepository) FindByID(ctx context.Context, id string) (*domain.Secret, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var secret domain.Secret
	err = r.db.Collection(secretCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&secret)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &secret, err
}

// FindByKey finds a secret by key, tenant, and environment
func (r *SecretRepository) FindByKey(ctx context.Context, tenantID, environment, key string) (*domain.Secret, error) {
	filter := bson.M{
		"secret_key":  key,
		"environment": environment,
	}
	if tenantID != "" {
		filter["tenant_id"] = tenantID
	}

	var secret domain.Secret
	err := r.db.Collection(secretCollection).FindOne(ctx, filter).Decode(&secret)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &secret, err
}

// Update updates a secret
func (r *SecretRepository) Update(ctx context.Context, secret *domain.Secret) error {
	secret.UpdatedAt = time.Now()
	secret.Version++

	filter := bson.M{"_id": secret.ID}
	update := bson.M{"$set": secret}

	_, err := r.db.Collection(secretCollection).UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a secret
func (r *SecretRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(secretCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// List lists secrets with pagination and filters (values masked)
func (r *SecretRepository) List(ctx context.Context, tenantID, environment string, page, perPage int) ([]*domain.Secret, int64, error) {
	filter := bson.M{}
	if tenantID != "" {
		filter["tenant_id"] = tenantID
	}
	if environment != "" {
		filter["environment"] = environment
	}

	// Count total
	total, err := r.db.Collection(secretCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate pagination
	skip := int64((page - 1) * perPage)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(perPage)).
		SetSort(bson.M{"updated_at": -1})

	cursor, err := r.db.Collection(secretCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var secrets []*domain.Secret
	if err = cursor.All(ctx, &secrets); err != nil {
		return nil, 0, err
	}

	return secrets, total, nil
}

// IncrementAccessCount increments the access count for a secret
func (r *SecretRepository) IncrementAccessCount(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$inc": bson.M{"access_count": 1},
		"$set": bson.M{"last_accessed_at": now},
	}

	_, err = r.db.Collection(secretCollection).UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)
	return err
}

// CreateAccessLog creates a secret access log entry
func (r *SecretRepository) CreateAccessLog(ctx context.Context, log *domain.SecretAccessLog) error {
	log.Timestamp = time.Now()

	_, err := r.db.Collection(secretAccessLogCollection).InsertOne(ctx, log)
	return err
}

// GetAccessLogs gets access logs for a secret
func (r *SecretRepository) GetAccessLogs(ctx context.Context, secretID string, page, perPage int) ([]*domain.SecretAccessLog, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(secretID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{"secret_id": objectID}

	// Count total
	total, err := r.db.Collection(secretAccessLogCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate pagination
	skip := int64((page - 1) * perPage)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(perPage)).
		SetSort(bson.M{"timestamp": -1})

	cursor, err := r.db.Collection(secretAccessLogCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*domain.SecretAccessLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetSecretsNeedingRotation gets secrets that need rotation
func (r *SecretRepository) GetSecretsNeedingRotation(ctx context.Context) ([]*domain.Secret, error) {
	filter := bson.M{
		"rotation_policy": "auto",
		"status":          "active",
	}

	cursor, err := r.db.Collection(secretCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var secrets []*domain.Secret
	for cursor.Next(ctx) {
		var secret domain.Secret
		if err := cursor.Decode(&secret); err != nil {
			// Log error but continue processing other secrets
			continue
		}

		// Check if rotation is needed
		if secret.LastRotatedAt != nil && secret.RotationDays > 0 {
			daysSinceRotation := time.Since(*secret.LastRotatedAt).Hours() / 24
			if daysSinceRotation >= float64(secret.RotationDays) {
				secrets = append(secrets, &secret)
			}
		}
	}

	// Check for errors during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return secrets, nil
}
