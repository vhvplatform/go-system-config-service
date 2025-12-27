package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Config represents a configuration entry with versioning support
type Config struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	TenantID    string                 `json:"tenant_id" bson:"tenant_id"`     // Optional: for tenant-specific configs
	ConfigKey   string                 `json:"config_key" bson:"config_key"`   // Unique configuration key
	Value       interface{}            `json:"value" bson:"value"`             // Configuration value (can be any type)
	Environment string                 `json:"environment" bson:"environment"` // dev, staging, production
	Version     int                    `json:"version" bson:"version"`         // Current version number
	Status      string                 `json:"status" bson:"status"`           // active, inactive, archived
	Description string                 `json:"description" bson:"description"` // Description of the configuration
	Tags        []string               `json:"tags" bson:"tags"`               // Tags for categorization
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`       // Additional metadata
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
	CreatedBy   string                 `json:"created_by" bson:"created_by"`
	UpdatedBy   string                 `json:"updated_by" bson:"updated_by"`
}

// ConfigVersion represents a version of a configuration
type ConfigVersion struct {
	ID              primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	ConfigID        primitive.ObjectID     `json:"config_id" bson:"config_id"`               // Reference to parent config
	ConfigKey       string                 `json:"config_key" bson:"config_key"`             // Denormalized for easier querying
	TenantID        string                 `json:"tenant_id" bson:"tenant_id"`               // Denormalized
	Environment     string                 `json:"environment" bson:"environment"`           // Denormalized
	VersionNumber   int                    `json:"version_number" bson:"version_number"`     // Version number
	Value           interface{}            `json:"value" bson:"value"`                       // Configuration value at this version
	ChangeReason    string                 `json:"change_reason" bson:"change_reason"`       // Reason for the change
	Status          string                 `json:"status" bson:"status"`                     // draft, active, archived
	IsActive        bool                   `json:"is_active" bson:"is_active"`               // Is this the active version?
	ValidationError string                 `json:"validation_error" bson:"validation_error"` // Any validation error
	Metadata        map[string]interface{} `json:"metadata" bson:"metadata"`                 // Additional metadata
	CreatedAt       time.Time              `json:"created_at" bson:"created_at"`
	CreatedBy       string                 `json:"created_by" bson:"created_by"`
}

// AuditLog represents an audit log entry for configuration changes
type AuditLog struct {
	ID           primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	ResourceType string                 `json:"resource_type" bson:"resource_type"` // config, secret, etc.
	ResourceID   primitive.ObjectID     `json:"resource_id" bson:"resource_id"`     // ID of the resource
	ResourceKey  string                 `json:"resource_key" bson:"resource_key"`   // Key of the resource
	TenantID     string                 `json:"tenant_id" bson:"tenant_id"`
	Environment  string                 `json:"environment" bson:"environment"`
	Action       string                 `json:"action" bson:"action"`         // create, update, delete, activate, rollback
	OldValue     interface{}            `json:"old_value" bson:"old_value"`   // Previous value
	NewValue     interface{}            `json:"new_value" bson:"new_value"`   // New value
	UserID       string                 `json:"user_id" bson:"user_id"`       // User who performed the action
	IPAddress    string                 `json:"ip_address" bson:"ip_address"` // IP address of the user
	UserAgent    string                 `json:"user_agent" bson:"user_agent"` // User agent
	Details      map[string]interface{} `json:"details" bson:"details"`       // Additional details
	Timestamp    time.Time              `json:"timestamp" bson:"timestamp"`
}

// Validate validates the configuration data
func (c *Config) Validate() error {
	if c.ConfigKey == "" {
		return errors.New("config_key is required")
	}
	if c.Environment == "" {
		return errors.New("environment is required")
	}
	// Validate environment value
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.Environment] {
		return errors.New("environment must be one of: development, staging, production")
	}
	if c.Status == "" {
		c.Status = "active"
	}
	// Validate status value
	validStatuses := map[string]bool{"active": true, "inactive": true, "archived": true}
	if !validStatuses[c.Status] {
		return errors.New("status must be one of: active, inactive, archived")
	}
	return nil
}

// Validate validates the configuration version data
func (cv *ConfigVersion) Validate() error {
	if cv.ConfigID.IsZero() {
		return errors.New("config_id is required")
	}
	if cv.ConfigKey == "" {
		return errors.New("config_key is required")
	}
	if cv.VersionNumber < 1 {
		return errors.New("version_number must be positive")
	}
	if cv.Status == "" {
		cv.Status = "draft"
	}
	// Validate status value
	validStatuses := map[string]bool{"draft": true, "active": true, "archived": true}
	if !validStatuses[cv.Status] {
		return errors.New("status must be one of: draft, active, archived")
	}
	return nil
}

// Validate validates the audit log data
func (al *AuditLog) Validate() error {
	if al.ResourceType == "" {
		return errors.New("resource_type is required")
	}
	if al.ResourceID.IsZero() {
		return errors.New("resource_id is required")
	}
	if al.Action == "" {
		return errors.New("action is required")
	}
	return nil
}
