package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AppComponent represents an application component in the system
type AppComponent struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	TenantID    string                 `json:"tenant_id" bson:"tenantId"`
	Code        string                 `json:"code" bson:"code"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Icon        string                 `json:"icon" bson:"icon"`
	Version     string                 `json:"version" bson:"version"`
	Status      string                 `json:"status" bson:"status"` // active, inactive
	Config      map[string]interface{} `json:"config" bson:"config"`
	CreatedAt   time.Time              `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updatedAt"`
	CreatedBy   string                 `json:"created_by" bson:"createdBy"`
	UpdatedBy   string                 `json:"updated_by" bson:"updatedBy"`
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
