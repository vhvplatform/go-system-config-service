package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// District represents a district/county in the system
type District struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code         string             `json:"code" bson:"code"`
	Name         map[string]string  `json:"name" bson:"name"` // i18n
	ProvinceCode string             `json:"province_code" bson:"provinceCode"`
	Type         string             `json:"type" bson:"type"` // district, county, etc.
	Status       string             `json:"status" bson:"status"`
	CreatedAt    time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updatedAt"`
}
