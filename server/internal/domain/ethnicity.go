package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ethnicity represents an ethnicity in the system
type Ethnicity struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code        string             `json:"code" bson:"code"`
	Name        map[string]string  `json:"name" bson:"name"` // i18n
	CountryCode string             `json:"country_code" bson:"countryCode"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updatedAt"`
}
