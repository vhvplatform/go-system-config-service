package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AppComponent represents an application component in the system
type AppComponent struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	TenantID    string                 `json:"tenant_id" bson:"tenant_id"`
	Code        string                 `json:"code" bson:"code"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Icon        string                 `json:"icon" bson:"icon"`
	Version     string                 `json:"version" bson:"version"`
	Status      string                 `json:"status" bson:"status"` // active, inactive
	Config      map[string]interface{} `json:"config" bson:"config"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
	CreatedBy   string                 `json:"created_by" bson:"created_by"`
	UpdatedBy   string                 `json:"updated_by" bson:"updated_by"`
}

// Validate validates the app component data
func (a *AppComponent) Validate() error {
	if a.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if a.Code == "" {
		return errors.New("code is required")
	}
	return nil
}
