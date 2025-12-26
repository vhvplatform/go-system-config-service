package migrations

import (
	"context"
	"time"

	"github.com/vhvplatform/go-system-config-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedData seeds initial data into the database
func SeedData(db *mongo.Database) error {
	ctx := context.Background()

	// Seed countries
	if err := seedCountries(ctx, db); err != nil {
		return err
	}

	// Seed currencies
	if err := seedCurrencies(ctx, db); err != nil {
		return err
	}

	// Seed default roles
	if err := seedRoles(ctx, db); err != nil {
		return err
	}

	return nil
}

func seedCountries(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("countries")

	// Check if countries already exist
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Skip seeding if data already exists
	}

	countries := []interface{}{
		domain.Country{
			Code:       "VN",
			Code3:      "VNM",
			Name:       map[string]string{"en": "Vietnam", "vi": "Việt Nam"},
			NativeName: "Việt Nam",
			PhoneCode:  "+84",
			Currency:   "VND",
			Flag:       "🇻🇳",
			Region:     "Asia",
			Status:     "active",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		domain.Country{
			Code:       "US",
			Code3:      "USA",
			Name:       map[string]string{"en": "United States", "vi": "Hoa Kỳ"},
			NativeName: "United States",
			PhoneCode:  "+1",
			Currency:   "USD",
			Flag:       "🇺🇸",
			Region:     "Americas",
			Status:     "active",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		domain.Country{
			Code:       "GB",
			Code3:      "GBR",
			Name:       map[string]string{"en": "United Kingdom", "vi": "Vương Quốc Anh"},
			NativeName: "United Kingdom",
			PhoneCode:  "+44",
			Currency:   "GBP",
			Flag:       "🇬🇧",
			Region:     "Europe",
			Status:     "active",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	_, err = collection.InsertMany(ctx, countries)
	return err
}

func seedCurrencies(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("currencies")

	// Check if currencies already exist
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Skip seeding if data already exists
	}

	currencies := []interface{}{
		domain.Currency{
			Code:          "VND",
			Name:          map[string]string{"en": "Vietnamese Dong", "vi": "Đồng Việt Nam"},
			Symbol:        "₫",
			DecimalDigits: 0,
			Countries:     []string{"VN"},
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		domain.Currency{
			Code:          "USD",
			Name:          map[string]string{"en": "US Dollar", "vi": "Đô la Mỹ"},
			Symbol:        "$",
			DecimalDigits: 2,
			Countries:     []string{"US"},
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		domain.Currency{
			Code:          "EUR",
			Name:          map[string]string{"en": "Euro", "vi": "Euro"},
			Symbol:        "€",
			DecimalDigits: 2,
			Countries:     []string{"DE", "FR", "IT", "ES"},
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		domain.Currency{
			Code:          "GBP",
			Name:          map[string]string{"en": "British Pound", "vi": "Bảng Anh"},
			Symbol:        "£",
			DecimalDigits: 2,
			Countries:     []string{"GB"},
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	_, err = collection.InsertMany(ctx, currencies)
	return err
}

func seedRoles(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("roles")

	// Check if roles already exist
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Skip seeding if data already exists
	}

	// Note: These are system-level default roles without tenant_id
	// Tenant-specific roles should be created when a tenant is provisioned
	roles := []interface{}{
		domain.Role{
			TenantID:    "", // Empty for system-level default roles
			Code:        "super_admin",
			Name:        "Super Administrator",
			Description: "Full system access",
			IsSystem:    true,
			Level:       1,
			Permissions: []string{"*"},
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		domain.Role{
			TenantID:    "", // Empty for system-level default roles
			Code:        "admin",
			Name:        "Administrator",
			Description: "Admin access to manage system",
			IsSystem:    true,
			Level:       2,
			Permissions: []string{"admin.*"},
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		domain.Role{
			TenantID:    "", // Empty for system-level default roles
			Code:        "manager",
			Name:        "Manager",
			Description: "Manager access",
			IsSystem:    true,
			Level:       3,
			Permissions: []string{"users.read", "users.update"},
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		domain.Role{
			TenantID:    "", // Empty for system-level default roles
			Code:        "user",
			Name:        "User",
			Description: "Standard user access",
			IsSystem:    true,
			Level:       4,
			Permissions: []string{"profile.read", "profile.update"},
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	_, err = collection.InsertMany(ctx, roles)
	return err
}
