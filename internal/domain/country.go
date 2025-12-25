package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Country represents a country in the system
type Country struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code       string             `json:"code" bson:"code"`   // ISO 3166-1 alpha-2
	Code3      string             `json:"code3" bson:"code3"` // ISO 3166-1 alpha-3
	Name       map[string]string  `json:"name" bson:"name"`   // i18n: en, vi
	NativeName string             `json:"native_name" bson:"native_name"`
	PhoneCode  string             `json:"phone_code" bson:"phone_code"`
	Currency   string             `json:"currency" bson:"currency"`
	Flag       string             `json:"flag" bson:"flag"`     // emoji or URL
	Region     string             `json:"region" bson:"region"` // Asia, Europe, etc.
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// Validate validates the country data
func (c *Country) Validate() error {
	if c.Code == "" {
		return errors.New("code is required")
	}
	if len(c.Name) == 0 {
		return errors.New("name is required")
	}
	return nil
}
