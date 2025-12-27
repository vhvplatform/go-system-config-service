package domain

// Common constants used across the application

// Valid environment values
const (
	EnvironmentDevelopment = "development"
	EnvironmentStaging     = "staging"
	EnvironmentProduction  = "production"
)

// Valid status values for configurations
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusArchived = "archived"
)

// Valid status values for configuration versions
const (
	VersionStatusDraft    = "draft"
	VersionStatusActive   = "active"
	VersionStatusArchived = "archived"
)

// Valid status values for secrets
const (
	SecretStatusActive  = "active"
	SecretStatusExpired = "expired"
	SecretStatusRotated = "rotated"
)

// Valid rotation policies for secrets
const (
	RotationPolicyManual = "manual"
	RotationPolicyAuto   = "auto"
)

// Valid status values for watch subscriptions
const (
	WatchStatusActive   = "active"
	WatchStatusPaused   = "paused"
	WatchStatusInactive = "inactive"
)

// GetValidEnvironments returns all valid environment values
func GetValidEnvironments() []string {
	return []string{EnvironmentDevelopment, EnvironmentStaging, EnvironmentProduction}
}

// IsValidEnvironment checks if an environment string is valid
func IsValidEnvironment(env string) bool {
	validEnvs := map[string]bool{
		EnvironmentDevelopment: true,
		EnvironmentStaging:     true,
		EnvironmentProduction:  true,
	}
	return validEnvs[env]
}
