package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServicePackage represents a service package in the system
type ServicePackage struct {
	ID           primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	TenantID     string                 `json:"tenant_id" bson:"tenantId"`
	Code         string                 `json:"code" bson:"code"`
	Name         string                 `json:"name" bson:"name"`
	Description  string                 `json:"description" bson:"description"`
	Tier         string                 `json:"tier" bson:"tier"` // free, basic, professional, enterprise
	Price        float64                `json:"price" bson:"price"`
	Currency     string                 `json:"currency" bson:"currency"`
	BillingCycle string                 `json:"billing_cycle" bson:"billingCycle"` // monthly, yearly
	Modules      []string               `json:"modules" bson:"modules"`            // module codes
	Limits       map[string]interface{} `json:"limits" bson:"limits"`              // users, storage, api_calls
	Features     []string               `json:"features" bson:"features"`
	IsPopular    bool                   `json:"is_popular" bson:"isPopular"`
	Status       string                 `json:"status" bson:"status"`
	CreatedAt    time.Time              `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updatedAt"`
}
