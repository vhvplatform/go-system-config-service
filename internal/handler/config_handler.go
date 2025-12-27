package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-system-config-service/internal/domain"
	"github.com/vhvplatform/go-system-config-service/internal/service"
	"go.uber.org/zap"
)

// ConfigHandler handles configuration HTTP requests
type ConfigHandler struct {
	service *service.ConfigService
	logger  *logger.Logger
}

// NewConfigHandler creates a new configuration handler
func NewConfigHandler(service *service.ConfigService, log *logger.Logger) *ConfigHandler {
	return &ConfigHandler{
		service: service,
		logger:  log,
	}
}

// Create creates a new configuration
// @Summary Create configuration
// @Tags configs
// @Accept json
// @Produce json
// @Param config body domain.Config true "Configuration data"
// @Success 201 {object} domain.Config
// @Router /api/v1/configs [post]
func (h *ConfigHandler) Create(c *gin.Context) {
	var config domain.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (would come from auth middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system" // Default for testing
	}
	config.CreatedBy = userID

	if err := h.service.Create(c.Request.Context(), &config); err != nil {
		h.logger.Error("Failed to create config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// GetByID gets a configuration by ID
// @Summary Get configuration by ID
// @Tags configs
// @Produce json
// @Param id path string true "Configuration ID"
// @Success 200 {object} domain.Config
// @Router /api/v1/configs/{id} [get]
func (h *ConfigHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	config, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get config", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GetByKey gets a configuration by key
// @Summary Get configuration by key
// @Tags configs
// @Produce json
// @Param key path string true "Configuration key"
// @Param tenant_id query string false "Tenant ID"
// @Param environment query string true "Environment"
// @Success 200 {object} domain.Config
// @Router /api/v1/configs/key/{key} [get]
func (h *ConfigHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	tenantID := c.Query("tenant_id")
	environment := c.Query("environment")

	if environment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "environment is required"})
		return
	}

	config, err := h.service.GetByKey(c.Request.Context(), tenantID, environment, key)
	if err != nil {
		h.logger.Error("Failed to get config", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// Update updates a configuration
// @Summary Update configuration
// @Tags configs
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Param updates body map[string]interface{} true "Updates"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/configs/{id} [put]
func (h *ConfigHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	if err := h.service.Update(c.Request.Context(), id, updates, userID); err != nil {
		h.logger.Error("Failed to update config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

// Delete deletes a configuration
// @Summary Delete configuration
// @Tags configs
// @Param id path string true "Configuration ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/configs/{id} [delete]
func (h *ConfigHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		h.logger.Error("Failed to delete config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}

// List lists configurations with pagination
// @Summary List configurations
// @Tags configs
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param environment query string false "Environment"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(30)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/configs [get]
func (h *ConfigHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	environment := c.Query("environment")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	configs, total, err := h.service.List(c.Request.Context(), tenantID, environment, page, perPage)
	if err != nil {
		h.logger.Error("Failed to list configs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list configurations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     configs,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// GetHistory gets version history for a configuration
// @Summary Get configuration history
// @Tags configs
// @Produce json
// @Param id path string true "Configuration ID"
// @Success 200 {array} domain.ConfigVersion
// @Router /api/v1/configs/{id}/history [get]
func (h *ConfigHandler) GetHistory(c *gin.Context) {
	id := c.Param("id")

	versions, err := h.service.GetHistory(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get version history"})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// ActivateVersion activates a specific version
// @Summary Activate configuration version
// @Tags configs
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Param body body map[string]interface{} true "Version info"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/configs/{id}/activate [post]
func (h *ConfigHandler) ActivateVersion(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		VersionNumber int `json:"version_number"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	if err := h.service.ActivateVersion(c.Request.Context(), id, body.VersionNumber, userID); err != nil {
		h.logger.Error("Failed to activate version", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Version activated successfully"})
}

// Rollback rolls back to a previous version
// @Summary Rollback configuration
// @Tags configs
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Param body body map[string]interface{} true "Target version"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/configs/{id}/rollback [post]
func (h *ConfigHandler) Rollback(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		TargetVersion int `json:"target_version"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	if err := h.service.Rollback(c.Request.Context(), id, body.TargetVersion, userID); err != nil {
		h.logger.Error("Failed to rollback", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration rolled back successfully"})
}

// CompareVersions compares two versions
// @Summary Compare configuration versions
// @Tags configs
// @Produce json
// @Param id path string true "Configuration ID"
// @Param v1 query int true "Version 1"
// @Param v2 query int true "Version 2"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/configs/{id}/compare [get]
func (h *ConfigHandler) CompareVersions(c *gin.Context) {
	id := c.Param("id")
	v1, _ := strconv.Atoi(c.Query("v1"))
	v2, _ := strconv.Atoi(c.Query("v2"))

	if v1 == 0 || v2 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "v1 and v2 parameters are required"})
		return
	}

	comparison, err := h.service.CompareVersions(c.Request.Context(), id, v1, v2)
	if err != nil {
		h.logger.Error("Failed to compare versions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compare versions"})
		return
	}

	c.JSON(http.StatusOK, comparison)
}

// GetAuditLogs gets audit logs for a configuration
// @Summary Get configuration audit logs
// @Tags configs
// @Produce json
// @Param id path string true "Configuration ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(30)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/configs/{id}/audit [get]
func (h *ConfigHandler) GetAuditLogs(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	logs, total, err := h.service.GetAuditLogs(c.Request.Context(), id, page, perPage)
	if err != nil {
		h.logger.Error("Failed to get audit logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     logs,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}
