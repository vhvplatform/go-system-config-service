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

// SecretHandler handles secret HTTP requests
type SecretHandler struct {
	service *service.SecretService
	logger  *logger.Logger
}

// NewSecretHandler creates a new secret handler
func NewSecretHandler(service *service.SecretService, log *logger.Logger) *SecretHandler {
	return &SecretHandler{
		service: service,
		logger:  log,
	}
}

// Create creates a new secret
// @Summary Create secret
// @Tags secrets
// @Accept json
// @Produce json
// @Param secret body map[string]interface{} true "Secret data with value"
// @Success 201 {object} domain.Secret
// @Router /api/v1/secrets [post]
func (h *SecretHandler) Create(c *gin.Context) {
	var input struct {
		Secret domain.Secret `json:"secret"`
		Value  string        `json:"value"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	if err := h.service.Create(c.Request.Context(), &input.Secret, input.Value, userID); err != nil {
		h.logger.Error("Failed to create secret", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create secret"})
		return
	}

	// Clear encrypted value before returning
	input.Secret.EncryptedValue = ""
	c.JSON(http.StatusCreated, input.Secret)
}

// GetByKey gets a secret by key and returns the decrypted value
// @Summary Get secret by key
// @Tags secrets
// @Produce json
// @Param key path string true "Secret key"
// @Param tenant_id query string false "Tenant ID"
// @Param environment query string true "Environment"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/secrets/key/{key} [get]
func (h *SecretHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	tenantID := c.Query("tenant_id")
	environment := c.Query("environment")

	if environment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "environment is required"})
		return
	}

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	value, err := h.service.GetByKey(c.Request.Context(), tenantID, environment, key, userID)
	if err != nil {
		h.logger.Error("Failed to get secret", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

// Update updates a secret
// @Summary Update secret
// @Tags secrets
// @Accept json
// @Produce json
// @Param id path string true "Secret ID"
// @Param body body map[string]interface{} true "New value"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/secrets/{id} [put]
func (h *SecretHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Value string `json:"value"`
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

	if err := h.service.Update(c.Request.Context(), id, body.Value, userID); err != nil {
		h.logger.Error("Failed to update secret", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Secret updated successfully"})
}

// Rotate rotates a secret
// @Summary Rotate secret
// @Tags secrets
// @Accept json
// @Produce json
// @Param id path string true "Secret ID"
// @Param body body map[string]interface{} true "New value"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/secrets/{id}/rotate [post]
func (h *SecretHandler) Rotate(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Value string `json:"value"`
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

	if err := h.service.Rotate(c.Request.Context(), id, body.Value, userID); err != nil {
		h.logger.Error("Failed to rotate secret", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rotate secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Secret rotated successfully"})
}

// Delete deletes a secret
// @Summary Delete secret
// @Tags secrets
// @Param id path string true "Secret ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/secrets/{id} [delete]
func (h *SecretHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		h.logger.Error("Failed to delete secret", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Secret deleted successfully"})
}

// List lists secrets with masked values
// @Summary List secrets
// @Tags secrets
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param environment query string false "Environment"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(30)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/secrets [get]
func (h *SecretHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	environment := c.Query("environment")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	secrets, total, err := h.service.List(c.Request.Context(), tenantID, environment, page, perPage)
	if err != nil {
		h.logger.Error("Failed to list secrets", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list secrets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     secrets,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// GetAuditLogs gets audit logs for a secret
// @Summary Get secret audit logs
// @Tags secrets
// @Produce json
// @Param id path string true "Secret ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(30)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/secrets/{id}/audit [get]
func (h *SecretHandler) GetAuditLogs(c *gin.Context) {
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
