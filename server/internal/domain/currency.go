package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Currency represents a currency in the system
type Currency struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code          string             `json:"code" bson:"code"` // ISO 4217
	Name          map[string]string  `json:"name" bson:"name"` // i18n
	Symbol        string             `json:"symbol" bson:"symbol"`
	DecimalDigits int                `json:"decimal_digits" bson:"decimalDigits"`
	Countries     []string           `json:"countries" bson:"countries"` // country codes
	Status        string             `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updatedAt"`
}
