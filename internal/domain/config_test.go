package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid config",
			config: Config{
				ConfigKey:   "db.timeout",
				Environment: "production",
				Status:      "active",
			},
			wantErr: false,
		},
		{
			name: "Missing config_key",
			config: Config{
				Environment: "production",
			},
			wantErr: true,
			errMsg:  "config_key is required",
		},
		{
			name: "Missing environment",
			config: Config{
				ConfigKey: "db.timeout",
			},
			wantErr: true,
			errMsg:  "environment is required",
		},
		{
			name: "Invalid environment",
			config: Config{
				ConfigKey:   "db.timeout",
				Environment: "invalid",
			},
			wantErr: true,
			errMsg:  "environment must be one of: development, staging, production",
		},
		{
			name: "Invalid status",
			config: Config{
				ConfigKey:   "db.timeout",
				Environment: "production",
				Status:      "invalid",
			},
			wantErr: true,
			errMsg:  "status must be one of: active, inactive, archived",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		version ConfigVersion
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid version",
			version: ConfigVersion{
				ConfigID:      primitive.NewObjectID(),
				ConfigKey:     "db.timeout",
				VersionNumber: 1,
				Status:        "active",
			},
			wantErr: false,
		},
		{
			name: "Missing config_id",
			version: ConfigVersion{
				ConfigKey:     "db.timeout",
				VersionNumber: 1,
			},
			wantErr: true,
			errMsg:  "config_id is required",
		},
		{
			name: "Missing config_key",
			version: ConfigVersion{
				ConfigID:      primitive.NewObjectID(),
				VersionNumber: 1,
			},
			wantErr: true,
			errMsg:  "config_key is required",
		},
		{
			name: "Invalid version number",
			version: ConfigVersion{
				ConfigID:      primitive.NewObjectID(),
				ConfigKey:     "db.timeout",
				VersionNumber: 0,
			},
			wantErr: true,
			errMsg:  "version_number must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.version.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecret_Validate(t *testing.T) {
	tests := []struct {
		name    string
		secret  Secret
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid secret",
			secret: Secret{
				SecretKey:      "db_password",
				Environment:    "production",
				Status:         "active",
				RotationPolicy: "manual",
			},
			wantErr: false,
		},
		{
			name: "Missing secret_key",
			secret: Secret{
				Environment: "production",
			},
			wantErr: true,
			errMsg:  "secret_key is required",
		},
		{
			name: "Missing environment",
			secret: Secret{
				SecretKey: "db_password",
			},
			wantErr: true,
			errMsg:  "environment is required",
		},
		{
			name: "Invalid environment",
			secret: Secret{
				SecretKey:   "db_password",
				Environment: "invalid",
			},
			wantErr: true,
			errMsg:  "environment must be one of: development, staging, production",
		},
		{
			name: "Invalid rotation policy",
			secret: Secret{
				SecretKey:      "db_password",
				Environment:    "production",
				RotationPolicy: "invalid",
			},
			wantErr: true,
			errMsg:  "rotation_policy must be one of: manual, auto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.secret.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecret_MaskedValue(t *testing.T) {
	secret := Secret{
		SecretKey:      "test_secret",
		EncryptedValue: "encrypted_data_here",
	}
	
	masked := secret.MaskedValue()
	assert.Equal(t, "***MASKED***", masked)
	assert.NotEqual(t, secret.EncryptedValue, masked)
}

func TestWatchSubscription_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sub     WatchSubscription
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid subscription",
			sub: WatchSubscription{
				SubscriberID: "service-1",
				CallbackURL:  "http://example.com/webhook",
				Patterns:     []string{"db.*", "api.*"},
				Status:       "active",
			},
			wantErr: false,
		},
		{
			name: "Missing subscriber_id",
			sub: WatchSubscription{
				CallbackURL: "http://example.com/webhook",
				Patterns:    []string{"db.*"},
			},
			wantErr: true,
			errMsg:  "subscriber_id is required",
		},
		{
			name: "Missing callback_url",
			sub: WatchSubscription{
				SubscriberID: "service-1",
				Patterns:     []string{"db.*"},
			},
			wantErr: true,
			errMsg:  "callback_url is required",
		},
		{
			name: "Missing patterns",
			sub: WatchSubscription{
				SubscriberID: "service-1",
				CallbackURL:  "http://example.com/webhook",
				Patterns:     []string{},
			},
			wantErr: true,
			errMsg:  "at least one pattern is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sub.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuditLog_Validate(t *testing.T) {
	tests := []struct {
		name    string
		log     AuditLog
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid audit log",
			log: AuditLog{
				ResourceType: "config",
				ResourceID:   primitive.NewObjectID(),
				Action:       "create",
			},
			wantErr: false,
		},
		{
			name: "Missing resource_type",
			log: AuditLog{
				ResourceID: primitive.NewObjectID(),
				Action:     "create",
			},
			wantErr: true,
			errMsg:  "resource_type is required",
		},
		{
			name: "Missing resource_id",
			log: AuditLog{
				ResourceType: "config",
				Action:       "create",
			},
			wantErr: true,
			errMsg:  "resource_id is required",
		},
		{
			name: "Missing action",
			log: AuditLog{
				ResourceType: "config",
				ResourceID:   primitive.NewObjectID(),
			},
			wantErr: true,
			errMsg:  "action is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
