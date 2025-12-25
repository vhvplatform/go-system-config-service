package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role represents a role in the system
type Role struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TenantID    string             `json:"tenant_id" bson:"tenant_id"`
	Code        string             `json:"code" bson:"code"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	IsSystem    bool               `json:"is_system" bson:"is_system"`     // system roles can't be deleted
	Level       int                `json:"level" bson:"level"`             // hierarchy level
	Permissions []string           `json:"permissions" bson:"permissions"` // permission codes
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
