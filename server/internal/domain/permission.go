package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Permission represents a permission in the system
type Permission struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TenantID    string             `json:"tenant_id" bson:"tenantId"`
	ModuleCode  string             `json:"module_code" bson:"moduleCode"`
	Code        string             `json:"code" bson:"code"` // users.create, users.read, users.update, users.delete
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Resource    string             `json:"resource" bson:"resource"` // users, products, orders
	Action      string             `json:"action" bson:"action"`     // create, read, update, delete, list, export
	Category    string             `json:"category" bson:"category"` // data, feature, admin
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updatedAt"`
}
