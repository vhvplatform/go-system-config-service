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
	configCollection        = "configs"
	configVersionCollection = "config_versions"
	auditLogCollection      = "config_audit_log"
)

// ConfigRepository handles configuration data persistence
type ConfigRepository struct {
	db *mongo.Database
}

// NewConfigRepository creates a new configuration repository
func NewConfigRepository(db *mongo.Database) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// Create creates a new configuration
func (r *ConfigRepository) Create(ctx context.Context, config *domain.Config) error {
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	config.Version = 1
	
	result, err := r.db.Collection(configCollection).InsertOne(ctx, config)
	if err != nil {
		return err
	}
	config.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a configuration by ID
func (r *ConfigRepository) FindByID(ctx context.Context, id string) (*domain.Config, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var config domain.Config
	err = r.db.Collection(configCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&config)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &config, err
}

// FindByKey finds a configuration by key, tenant, and environment
func (r *ConfigRepository) FindByKey(ctx context.Context, tenantID, environment, key string) (*domain.Config, error) {
	filter := bson.M{
		"config_key":  key,
		"environment": environment,
	}
	if tenantID != "" {
		filter["tenant_id"] = tenantID
	}
	
	var config domain.Config
	err := r.db.Collection(configCollection).FindOne(ctx, filter).Decode(&config)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &config, err
}

// Update updates a configuration
func (r *ConfigRepository) Update(ctx context.Context, config *domain.Config) error {
	config.UpdatedAt = time.Now()
	config.Version++
	
	filter := bson.M{"_id": config.ID}
	update := bson.M{"$set": config}
	
	_, err := r.db.Collection(configCollection).UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a configuration
func (r *ConfigRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Collection(configCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// List lists configurations with pagination and filters
func (r *ConfigRepository) List(ctx context.Context, tenantID, environment string, page, perPage int) ([]*domain.Config, int64, error) {
	filter := bson.M{}
	if tenantID != "" {
		filter["tenant_id"] = tenantID
	}
	if environment != "" {
		filter["environment"] = environment
	}
	
	// Count total
	total, err := r.db.Collection(configCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	// Calculate pagination
	skip := int64((page - 1) * perPage)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(perPage)).
		SetSort(bson.M{"updated_at": -1})
	
	cursor, err := r.db.Collection(configCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var configs []*domain.Config
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, 0, err
	}
	
	return configs, total, nil
}

// CreateVersion creates a new configuration version
func (r *ConfigRepository) CreateVersion(ctx context.Context, version *domain.ConfigVersion) error {
	version.CreatedAt = time.Now()
	
	result, err := r.db.Collection(configVersionCollection).InsertOne(ctx, version)
	if err != nil {
		return err
	}
	version.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetVersionHistory gets all versions of a configuration
func (r *ConfigRepository) GetVersionHistory(ctx context.Context, configID string) ([]*domain.ConfigVersion, error) {
	objectID, err := primitive.ObjectIDFromHex(configID)
	if err != nil {
		return nil, err
	}
	
	opts := options.Find().SetSort(bson.M{"version_number": -1})
	cursor, err := r.db.Collection(configVersionCollection).Find(
		ctx,
		bson.M{"config_id": objectID},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var versions []*domain.ConfigVersion
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, err
	}
	
	return versions, nil
}

// GetVersion gets a specific version
func (r *ConfigRepository) GetVersion(ctx context.Context, configID string, versionNumber int) (*domain.ConfigVersion, error) {
	objectID, err := primitive.ObjectIDFromHex(configID)
	if err != nil {
		return nil, err
	}
	
	var version domain.ConfigVersion
	err = r.db.Collection(configVersionCollection).FindOne(
		ctx,
		bson.M{
			"config_id":      objectID,
			"version_number": versionNumber,
		},
	).Decode(&version)
	
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &version, err
}

// ActivateVersion activates a specific version
func (r *ConfigRepository) ActivateVersion(ctx context.Context, configID string, versionNumber int) error {
	objectID, err := primitive.ObjectIDFromHex(configID)
	if err != nil {
		return err
	}
	
	// Deactivate all versions first
	_, err = r.db.Collection(configVersionCollection).UpdateMany(
		ctx,
		bson.M{"config_id": objectID},
		bson.M{"$set": bson.M{"is_active": false}},
	)
	if err != nil {
		return err
	}
	
	// Activate the specific version
	_, err = r.db.Collection(configVersionCollection).UpdateOne(
		ctx,
		bson.M{
			"config_id":      objectID,
			"version_number": versionNumber,
		},
		bson.M{
			"$set": bson.M{
				"is_active": true,
				"status":    "active",
			},
		},
	)
	return err
}

// CreateAuditLog creates an audit log entry
func (r *ConfigRepository) CreateAuditLog(ctx context.Context, log *domain.AuditLog) error {
	log.Timestamp = time.Now()
	
	_, err := r.db.Collection(auditLogCollection).InsertOne(ctx, log)
	return err
}

// GetAuditLogs gets audit logs with filters
func (r *ConfigRepository) GetAuditLogs(ctx context.Context, resourceID string, page, perPage int) ([]*domain.AuditLog, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {
		return nil, 0, err
	}
	
	filter := bson.M{"resource_id": objectID}
	
	// Count total
	total, err := r.db.Collection(auditLogCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	// Calculate pagination
	skip := int64((page - 1) * perPage)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(perPage)).
		SetSort(bson.M{"timestamp": -1})
	
	cursor, err := r.db.Collection(auditLogCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var logs []*domain.AuditLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}
	
	return logs, total, nil
}
