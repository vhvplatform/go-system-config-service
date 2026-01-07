package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SaaSModule represents a SaaS module in the system
type SaaSModule struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TenantID     string             `json:"tenant_id" bson:"tenantId"`
	Code         string             `json:"code" bson:"code"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	Icon         string             `json:"icon" bson:"icon"`
	Category     string             `json:"category" bson:"category"` // core, addon, premium
	IsCore       bool               `json:"is_core" bson:"isCore"`
	Dependencies []string           `json:"dependencies" bson:"dependencies"` // module codes
	Price        float64            `json:"price" bson:"price"`
	Status       string             `json:"status" bson:"status"` // active, inactive, deprecated
	Features     []string           `json:"features" bson:"features"`
	CreatedAt    time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updatedAt"`
}
