package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminMenu represents an admin menu item in the system
type AdminMenu struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TenantID    string             `json:"tenant_id" bson:"tenantId"`
	ModuleCode  string             `json:"module_code" bson:"moduleCode"`
	ParentID    string             `json:"parent_id" bson:"parentId"`
	Code        string             `json:"code" bson:"code"`
	Name        string             `json:"name" bson:"name"`
	Title       map[string]string  `json:"title" bson:"title"` // i18n: en, vi
	Icon        string             `json:"icon" bson:"icon"`
	Path        string             `json:"path" bson:"path"`
	Component   string             `json:"component" bson:"component"`
	Order       int                `json:"order" bson:"order"`
	Permissions []string           `json:"permissions" bson:"permissions"` // required permissions
	IsVisible   bool               `json:"is_visible" bson:"isVisible"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updatedAt"`
}
