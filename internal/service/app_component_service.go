package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vhvcorp/go-shared/errors"
	"github.com/vhvcorp/go-shared/logger"
	"github.com/vhvcorp/go-shared/redis"
	"github.com/vhvcorp/go-system-config-service/internal/domain"
	"github.com/vhvcorp/go-system-config-service/internal/repository"
	"go.uber.org/zap"
)

// AppComponentService handles app component business logic
type AppComponentService struct {
	repo        *repository.AppComponentRepository
	redisClient *redis.Client
	logger      *logger.Logger
}

// NewAppComponentService creates a new app component service
func NewAppComponentService(
	repo *repository.AppComponentRepository,
	redisClient *redis.Client,
	log *logger.Logger,
) *AppComponentService {
	return &AppComponentService{
		repo:        repo,
		redisClient: redisClient,
		logger:      log,
	}
}

// Create creates a new app component
func (s *AppComponentService) Create(ctx context.Context, component *domain.AppComponent) error {
	// Check if component with same code exists
	existing, err := s.repo.FindByCode(ctx, component.TenantID, component.Code)
	if err != nil {
		s.logger.Error("Failed to check existing component", zap.Error(err))
		return errors.Internal("Failed to create app component")
	}
	if existing != nil {
		return errors.Conflict("App component with this code already exists")
	}

	if err := s.repo.Create(ctx, component); err != nil {
		s.logger.Error("Failed to create app component", zap.Error(err))
		return errors.Internal("Failed to create app component")
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("system-config:app-component:%s:%s", component.TenantID, component.Code)
	s.redisClient.Delete(ctx, cacheKey)

	s.logger.Info("App component created", zap.String("id", component.ID.Hex()), zap.String("code", component.Code))
	return nil
}

// GetByID gets an app component by ID
func (s *AppComponentService) GetByID(ctx context.Context, id string) (*domain.AppComponent, error) {
	component, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get app component", zap.Error(err))
		return nil, errors.Internal("Failed to get app component")
	}
	if component == nil {
		return nil, errors.NotFound("App component not found")
	}
	return component, nil
}

// GetByCode gets an app component by code with caching
func (s *AppComponentService) GetByCode(ctx context.Context, tenantID, code string) (*domain.AppComponent, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("system-config:app-component:%s:%s", tenantID, code)
	cached, err := s.redisClient.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var component domain.AppComponent
		if err := json.Unmarshal([]byte(cached), &component); err == nil {
			return &component, nil
		}
	}

	// Get from database
	component, err := s.repo.FindByCode(ctx, tenantID, code)
	if err != nil {
		s.logger.Error("Failed to get app component", zap.Error(err))
		return nil, errors.Internal("Failed to get app component")
	}
	if component == nil {
		return nil, errors.NotFound("App component not found")
	}

	// Cache for 1 hour
	if data, err := json.Marshal(component); err == nil {
		s.redisClient.Set(ctx, cacheKey, data, 1*time.Hour)
	}

	return component, nil
}

// List lists app components with pagination
func (s *AppComponentService) List(ctx context.Context, tenantID string, page, perPage int) ([]*domain.AppComponent, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 30
	}

	components, total, err := s.repo.List(ctx, tenantID, page, perPage)
	if err != nil {
		s.logger.Error("Failed to list app components", zap.Error(err))
		return nil, 0, errors.Internal("Failed to list app components")
	}

	return components, total, nil
}

// Update updates an app component
func (s *AppComponentService) Update(ctx context.Context, component *domain.AppComponent) error {
	// Check if exists
	existing, err := s.repo.FindByID(ctx, component.ID.Hex())
	if err != nil {
		s.logger.Error("Failed to check existing component", zap.Error(err))
		return errors.Internal("Failed to update app component")
	}
	if existing == nil {
		return errors.NotFound("App component not found")
	}

	if err := s.repo.Update(ctx, component); err != nil {
		s.logger.Error("Failed to update app component", zap.Error(err))
		return errors.Internal("Failed to update app component")
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("system-config:app-component:%s:%s", component.TenantID, component.Code)
	s.redisClient.Delete(ctx, cacheKey)

	s.logger.Info("App component updated", zap.String("id", component.ID.Hex()))
	return nil
}

// Delete deletes an app component
func (s *AppComponentService) Delete(ctx context.Context, id, tenantID string) error {
	// Check if exists
	component, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to check existing component", zap.Error(err))
		return errors.Internal("Failed to delete app component")
	}
	if component == nil {
		return errors.NotFound("App component not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete app component", zap.Error(err))
		return errors.Internal("Failed to delete app component")
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("system-config:app-component:%s:%s", component.TenantID, component.Code)
	s.redisClient.Delete(ctx, cacheKey)

	s.logger.Info("App component deleted", zap.String("id", id))
	return nil
}
