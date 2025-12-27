package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Secret represents a sensitive configuration with encryption
type Secret struct {
	ID              primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	TenantID        string                 `json:"tenant_id" bson:"tenant_id"`                 // Optional: for tenant-specific secrets
	SecretKey       string                 `json:"secret_key" bson:"secret_key"`               // Unique secret key
	EncryptedValue  string                 `json:"-" bson:"encrypted_value"`                   // Encrypted value (not exposed in JSON)
	Environment     string                 `json:"environment" bson:"environment"`             // dev, staging, production
	Description     string                 `json:"description" bson:"description"`             // Description of the secret
	RotationPolicy  string                 `json:"rotation_policy" bson:"rotation_policy"`     // manual, auto
	RotationDays    int                    `json:"rotation_days" bson:"rotation_days"`         // Days until rotation required
	LastRotatedAt   *time.Time             `json:"last_rotated_at" bson:"last_rotated_at"`     // Last rotation timestamp
	ExpiresAt       *time.Time             `json:"expires_at" bson:"expires_at"`               // Optional expiration
	Status          string                 `json:"status" bson:"status"`                       // active, expired, rotated
	Version         int                    `json:"version" bson:"version"`                     // Secret version
	EncryptionKeyID string                 `json:"encryption_key_id" bson:"encryption_key_id"` // ID of encryption key used
	Metadata        map[string]interface{} `json:"metadata" bson:"metadata"`                   // Additional metadata
	AccessCount     int64                  `json:"access_count" bson:"access_count"`           // Number of times accessed
	LastAccessedAt  *time.Time             `json:"last_accessed_at" bson:"last_accessed_at"`   // Last access timestamp
	CreatedAt       time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" bson:"updated_at"`
	CreatedBy       string                 `json:"created_by" bson:"created_by"`
	UpdatedBy       string                 `json:"updated_by" bson:"updated_by"`
}

// SecretAccessLog represents an access log for secrets
type SecretAccessLog struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	SecretID    primitive.ObjectID     `json:"secret_id" bson:"secret_id"`       // Reference to secret
	SecretKey   string                 `json:"secret_key" bson:"secret_key"`     // Denormalized for easier querying
	TenantID    string                 `json:"tenant_id" bson:"tenant_id"`       // Denormalized
	Environment string                 `json:"environment" bson:"environment"`   // Denormalized
	UserID      string                 `json:"user_id" bson:"user_id"`           // User who accessed the secret
	ServiceName string                 `json:"service_name" bson:"service_name"` // Service that accessed
	Action      string                 `json:"action" bson:"action"`             // read, update, rotate, delete
	IPAddress   string                 `json:"ip_address" bson:"ip_address"`     // IP address
	UserAgent   string                 `json:"user_agent" bson:"user_agent"`     // User agent
	Success     bool                   `json:"success" bson:"success"`           // Was access successful
	FailReason  string                 `json:"fail_reason" bson:"fail_reason"`   // Reason if failed
	Details     map[string]interface{} `json:"details" bson:"details"`           // Additional details
	Timestamp   time.Time              `json:"timestamp" bson:"timestamp"`
}

// WatchSubscription represents a subscription to configuration changes
type WatchSubscription struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SubscriberID string             `json:"subscriber_id" bson:"subscriber_id"` // Unique ID for subscriber
	TenantID     string             `json:"tenant_id" bson:"tenant_id"`         // Optional: tenant filter
	ServiceName  string             `json:"service_name" bson:"service_name"`   // Name of subscribing service
	CallbackURL  string             `json:"callback_url" bson:"callback_url"`   // Webhook URL for notifications
	Patterns     []string           `json:"patterns" bson:"patterns"`           // Config key patterns to watch (e.g., "db.*", "api.*.timeout")
	Environments []string           `json:"environments" bson:"environments"`   // Environments to watch
	Status       string             `json:"status" bson:"status"`               // active, paused, inactive
	LastNotified *time.Time         `json:"last_notified" bson:"last_notified"` // Last notification time
	FailureCount int                `json:"failure_count" bson:"failure_count"` // Number of consecutive failures
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// ConfigChangeNotification represents a notification for configuration changes
type ConfigChangeNotification struct {
	ConfigKey   string                 `json:"config_key"`
	TenantID    string                 `json:"tenant_id"`
	Environment string                 `json:"environment"`
	OldValue    interface{}            `json:"old_value"`
	NewValue    interface{}            `json:"new_value"`
	Version     int                    `json:"version"`
	ChangeType  string                 `json:"change_type"` // create, update, delete, activate, rollback
	ChangedBy   string                 `json:"changed_by"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Validate validates the secret data
func (s *Secret) Validate() error {
	if s.SecretKey == "" {
		return errors.New("secret_key is required")
	}
	if s.Environment == "" {
		return errors.New("environment is required")
	}
	// Validate environment value
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[s.Environment] {
		return errors.New("environment must be one of: development, staging, production")
	}
	if s.Status == "" {
		s.Status = "active"
	}
	// Validate status value
	validStatuses := map[string]bool{"active": true, "expired": true, "rotated": true}
	if !validStatuses[s.Status] {
		return errors.New("status must be one of: active, expired, rotated")
	}
	if s.RotationPolicy != "" {
		validPolicies := map[string]bool{"manual": true, "auto": true}
		if !validPolicies[s.RotationPolicy] {
			return errors.New("rotation_policy must be one of: manual, auto")
		}
	}
	return nil
}

// Validate validates the watch subscription data
func (w *WatchSubscription) Validate() error {
	if w.SubscriberID == "" {
		return errors.New("subscriber_id is required")
	}
	if w.CallbackURL == "" {
		return errors.New("callback_url is required")
	}
	if len(w.Patterns) == 0 {
		return errors.New("at least one pattern is required")
	}
	if w.Status == "" {
		w.Status = "active"
	}
	// Validate status value
	validStatuses := map[string]bool{"active": true, "paused": true, "inactive": true}
	if !validStatuses[w.Status] {
		return errors.New("status must be one of: active, paused, inactive")
	}
	return nil
}

// MaskedValue returns a masked version of the secret for display
func (s *Secret) MaskedValue() string {
	return "***MASKED***"
}
