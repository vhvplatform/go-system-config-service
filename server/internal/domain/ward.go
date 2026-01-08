package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ward represents a ward/commune in the system
type Ward struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code         string             `json:"code" bson:"code"`
	Name         map[string]string  `json:"name" bson:"name"` // i18n
	DistrictCode string             `json:"district_code" bson:"districtCode"`
	Type         string             `json:"type" bson:"type"` // ward, commune, etc.
	Status       string             `json:"status" bson:"status"`
	CreatedAt    time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updatedAt"`
}
