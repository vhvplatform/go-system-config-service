package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vhvplatform/go-shared/errors"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-shared/redis"
	"github.com/vhvplatform/go-system-config-service/internal/domain"
	"github.com/vhvplatform/go-system-config-service/internal/repository"
	"go.uber.org/zap"
)

// CountryService handles country business logic
type CountryService struct {
	repo        *repository.CountryRepository
	redisClient *redis.Client
	logger      *logger.Logger
}

// NewCountryService creates a new country service
func NewCountryService(
	repo *repository.CountryRepository,
	redisClient *redis.Client,
	log *logger.Logger,
) *CountryService {
	return &CountryService{
		repo:        repo,
		redisClient: redisClient,
		logger:      log,
	}
}

// Create creates a new country
func (s *CountryService) Create(ctx context.Context, country *domain.Country) error {
	// Check if country with same code exists
	existing, err := s.repo.FindByCode(ctx, country.Code)
	if err != nil {
		s.logger.Error("Failed to check existing country", zap.Error(err))
		return errors.Internal("Failed to create country")
	}
	if existing != nil {
		return errors.Conflict("Country with this code already exists")
	}

	if err := s.repo.Create(ctx, country); err != nil {
		s.logger.Error("Failed to create country", zap.Error(err))
		return errors.Internal("Failed to create country")
	}

	// Invalidate cache - use a pattern-based approach for better cache management
	cacheKey := fmt.Sprintf("system-config:country:%s", country.Code)
	// Invalidate list cache as well
	s.redisClient.Delete(ctx, cacheKey)
	s.redisClient.Delete(ctx, "system-config:countries:list:p1:30")

	s.logger.Info("Country created", zap.String("code", country.Code))
	return nil
}

// GetByCode gets a country by code with caching
func (s *CountryService) GetByCode(ctx context.Context, code string) (*domain.Country, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("system-config:country:%s", code)
	cached, err := s.redisClient.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		// Check for negative cache (non-existent record marker)
		if cached == "NOT_FOUND" {
			return nil, errors.NotFound("Country not found")
		}

		var country domain.Country
		if err := json.Unmarshal([]byte(cached), &country); err == nil {
			return &country, nil
		}
	}

	// Get from database
	country, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		s.logger.Error("Failed to get country", zap.Error(err))
		return nil, errors.Internal("Failed to get country")
	}

	if country == nil {
		// Implement negative caching: cache the fact that this country doesn't exist
		// This prevents repeated database hits for non-existent countries
		s.redisClient.Set(ctx, cacheKey, []byte("NOT_FOUND"), 5*time.Minute)
		return nil, errors.NotFound("Country not found")
	}

	// Cache for 24 hours (master data changes infrequently)
	if data, err := json.Marshal(country); err == nil {
		s.redisClient.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return country, nil
}

// List lists all countries with caching
func (s *CountryService) List(ctx context.Context, page, perPage int) ([]*domain.Country, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 30
	}

	// Try cache first for the first page with default page size (most common query)
	if page == 1 && perPage == 30 {
		cacheKey := "system-config:countries:list:p1:30"
		cached, err := s.redisClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var cachedData struct {
				Countries []*domain.Country `json:"countries"`
				Total     int64             `json:"total"`
			}
			if err := json.Unmarshal([]byte(cached), &cachedData); err == nil {
				return cachedData.Countries, cachedData.Total, nil
			}
		}

		// Get from database
		countries, total, err := s.repo.List(ctx, page, perPage)
		if err != nil {
			s.logger.Error("Failed to list countries", zap.Error(err))
			return nil, 0, errors.Internal("Failed to list countries")
		}

		// Cache the first page for 1 hour (frequently accessed)
		cachedData := struct {
			Countries []*domain.Country `json:"countries"`
			Total     int64             `json:"total"`
		}{
			Countries: countries,
			Total:     total,
		}
		if data, err := json.Marshal(cachedData); err == nil {
			s.redisClient.Set(ctx, cacheKey, data, 1*time.Hour)
		}

		return countries, total, nil
	}

	// For other pages, get directly from database (less frequently accessed)
	countries, total, err := s.repo.List(ctx, page, perPage)
	if err != nil {
		s.logger.Error("Failed to list countries", zap.Error(err))
		return nil, 0, errors.Internal("Failed to list countries")
	}

	return countries, total, nil
}

// Update updates a country
func (s *CountryService) Update(ctx context.Context, country *domain.Country) error {
	// Check if exists
	existing, err := s.repo.FindByCode(ctx, country.Code)
	if err != nil {
		s.logger.Error("Failed to check existing country", zap.Error(err))
		return errors.Internal("Failed to update country")
	}
	if existing == nil {
		return errors.NotFound("Country not found")
	}

	if err := s.repo.Update(ctx, country); err != nil {
		s.logger.Error("Failed to update country", zap.Error(err))
		return errors.Internal("Failed to update country")
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("system-config:country:%s", country.Code)
	s.redisClient.Delete(ctx, cacheKey)
	s.redisClient.Delete(ctx, "system-config:countries:list:p1:30")

	s.logger.Info("Country updated", zap.String("code", country.Code))
	return nil
}

// Delete deletes a country
func (s *CountryService) Delete(ctx context.Context, code string) error {
	// Check if exists
	country, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		s.logger.Error("Failed to check existing country", zap.Error(err))
		return errors.Internal("Failed to delete country")
	}
	if country == nil {
		return errors.NotFound("Country not found")
	}

	if err := s.repo.Delete(ctx, code); err != nil {
		s.logger.Error("Failed to delete country", zap.Error(err))
		return errors.Internal("Failed to delete country")
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("system-config:country:%s", code)
	s.redisClient.Delete(ctx, cacheKey)
	s.redisClient.Delete(ctx, "system-config:countries:list:p1:30")

	s.logger.Info("Country deleted", zap.String("code", code))
	return nil
}

// GetByCodes gets multiple countries by codes efficiently (batch operation)
// This method uses the new batch repository method to fetch multiple countries in one query
func (s *CountryService) GetByCodes(ctx context.Context, codes []string) ([]*domain.Country, error) {
	if len(codes) == 0 {
		return []*domain.Country{}, nil
	}

	// Try to get all from cache first
	cacheKeys := make([]string, len(codes))
	for i, code := range codes {
		cacheKeys[i] = fmt.Sprintf("system-config:country:%s", code)
	}

	// Get cached results
	cachedCountries := make(map[string]*domain.Country)
	missingCodes := []string{}

	for i, code := range codes {
		cached, err := s.redisClient.Get(ctx, cacheKeys[i])
		if err == nil && cached != "" && cached != "NOT_FOUND" {
			var country domain.Country
			if err := json.Unmarshal([]byte(cached), &country); err == nil {
				cachedCountries[code] = &country
				continue
			}
		}
		// Add to missing list if:
		// 1. Cache miss (cached is empty)
		// 2. Deserialization failed
		// Skip if explicitly marked as NOT_FOUND in cache (negative cache)
		if cached != "NOT_FOUND" {
			missingCodes = append(missingCodes, code)
		}
	}

	// Fetch missing countries from database in batch
	var dbCountries []*domain.Country
	if len(missingCodes) > 0 {
		var err error
		dbCountries, err = s.repo.FindByCodes(ctx, missingCodes)
		if err != nil {
			s.logger.Error("Failed to get countries", zap.Error(err))
			return nil, errors.Internal("Failed to get countries")
		}

		// Cache the retrieved countries
		for _, country := range dbCountries {
			if data, err := json.Marshal(country); err == nil {
				cacheKey := fmt.Sprintf("system-config:country:%s", country.Code)
				s.redisClient.Set(ctx, cacheKey, data, 24*time.Hour)
			}
			cachedCountries[country.Code] = country
		}

		// Implement negative caching for codes that weren't found
		foundCodes := make(map[string]bool)
		for _, country := range dbCountries {
			foundCodes[country.Code] = true
		}
		for _, code := range missingCodes {
			if !foundCodes[code] {
				cacheKey := fmt.Sprintf("system-config:country:%s", code)
				s.redisClient.Set(ctx, cacheKey, []byte("NOT_FOUND"), 5*time.Minute)
			}
		}
	}

	// Build result in the same order as requested codes
	result := make([]*domain.Country, 0, len(codes))
	for _, code := range codes {
		if country, exists := cachedCountries[code]; exists {
			result = append(result, country)
		}
	}

	return result, nil
}
