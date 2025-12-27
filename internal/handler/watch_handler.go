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

// WatchHandler handles watch subscription HTTP requests
type WatchHandler struct {
	service *service.WatchService
	logger  *logger.Logger
}

// NewWatchHandler creates a new watch handler
func NewWatchHandler(service *service.WatchService, log *logger.Logger) *WatchHandler {
	return &WatchHandler{
		service: service,
		logger:  log,
	}
}

// Subscribe creates a new watch subscription
// @Summary Subscribe to config changes
// @Tags watch
// @Accept json
// @Produce json
// @Param subscription body domain.WatchSubscription true "Subscription data"
// @Success 201 {object} domain.WatchSubscription
// @Router /api/v1/watch/subscribe [post]
func (h *WatchHandler) Subscribe(c *gin.Context) {
	var subscription domain.WatchSubscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Subscribe(c.Request.Context(), &subscription); err != nil {
		h.logger.Error("Failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// Unsubscribe removes a watch subscription
// @Summary Unsubscribe from config changes
// @Tags watch
// @Param id path string true "Subscription ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/watch/unsubscribe/{id} [delete]
func (h *WatchHandler) Unsubscribe(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Unsubscribe(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to unsubscribe", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

// List lists all subscriptions
// @Summary List watch subscriptions
// @Tags watch
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(30)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/watch/subscriptions [get]
func (h *WatchHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))

	subscriptions, total, err := h.service.List(c.Request.Context(), page, perPage)
	if err != nil {
		h.logger.Error("Failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     subscriptions,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// GetByID gets a subscription by ID
// @Summary Get watch subscription
// @Tags watch
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} domain.WatchSubscription
// @Router /api/v1/watch/subscriptions/{id} [get]
func (h *WatchHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	subscription, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get subscription", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// Update updates a subscription
// @Summary Update watch subscription
// @Tags watch
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param updates body map[string]interface{} true "Updates"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/watch/subscriptions/{id} [put]
func (h *WatchHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSubscription(c.Request.Context(), id, updates); err != nil {
		h.logger.Error("Failed to update subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
}

// TriggerNotification manually triggers a notification for testing
// @Summary Trigger manual notification
// @Tags watch
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Notification data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/watch/trigger [post]
func (h *WatchHandler) TriggerNotification(c *gin.Context) {
	var body struct {
		ConfigKey   string `json:"config_key"`
		TenantID    string `json:"tenant_id"`
		Environment string `json:"environment"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.TriggerNotification(c.Request.Context(), body.ConfigKey, body.TenantID, body.Environment); err != nil {
		h.logger.Error("Failed to trigger notification", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trigger notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification triggered successfully"})
}

// GetMatchingSubscriptions gets subscriptions that would match a given config key
// @Summary Get matching subscriptions
// @Tags watch
// @Produce json
// @Param config_key query string true "Configuration key"
// @Param tenant_id query string false "Tenant ID"
// @Param environment query string true "Environment"
// @Success 200 {array} domain.WatchSubscription
// @Router /api/v1/watch/matching [get]
func (h *WatchHandler) GetMatchingSubscriptions(c *gin.Context) {
	configKey := c.Query("config_key")
	tenantID := c.Query("tenant_id")
	environment := c.Query("environment")

	if configKey == "" || environment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config_key and environment are required"})
		return
	}

	subscriptions, err := h.service.GetMatchingSubscriptions(c.Request.Context(), configKey, tenantID, environment)
	if err != nil {
		h.logger.Error("Failed to get matching subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get matching subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}
