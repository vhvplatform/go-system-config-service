package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhvcorp/go-shared/errors"
	"github.com/vhvcorp/go-shared/logger"
	"github.com/vhvcorp/go-system-config-service/internal/domain"
	"github.com/vhvcorp/go-system-config-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// AppComponentHandler handles HTTP requests for app components
type AppComponentHandler struct {
	service *service.AppComponentService
	logger  *logger.Logger
}

// NewAppComponentHandler creates a new app component handler
func NewAppComponentHandler(service *service.AppComponentService, log *logger.Logger) *AppComponentHandler {
	return &AppComponentHandler{
		service: service,
		logger:  log,
	}
}

// Create handles creating a new app component
func (h *AppComponentHandler) Create(c *gin.Context) {
	var component domain.AppComponent
	if err := c.ShouldBindJSON(&component); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	// Get tenant ID from context (set by middleware)
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		h.respondError(c, errors.BadRequest("Tenant ID is required"))
		return
	}
	component.TenantID = tenantID

	if err := h.service.Create(c.Request.Context(), &component); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": component})
}

// GetByID handles getting an app component by ID
func (h *AppComponentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respondError(c, errors.BadRequest("ID is required"))
		return
	}

	component, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": component})
}

// List handles listing app components
func (h *AppComponentHandler) List(c *gin.Context) {
	var req domain.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PerPage = 30
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		h.respondError(c, errors.BadRequest("Tenant ID is required"))
		return
	}

	components, total, err := h.service.List(c.Request.Context(), tenantID, req.Page, req.PerPage)
	if err != nil {
		h.respondError(c, err)
		return
	}

	totalPages := int(total) / req.PerPage
	if int(total)%req.PerPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": components,
		"pagination": domain.PaginationResponse{
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: totalPages,
			TotalItems: total,
		},
	})
}

// Update handles updating an app component
func (h *AppComponentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respondError(c, errors.BadRequest("ID is required"))
		return
	}

	var component domain.AppComponent
	if err := c.ShouldBindJSON(&component); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	// Convert ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.respondError(c, errors.BadRequest("Invalid ID format"))
		return
	}
	component.ID = objectID

	if err := h.service.Update(c.Request.Context(), &component); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": component})
}

// Delete handles deleting an app component
func (h *AppComponentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respondError(c, errors.BadRequest("ID is required"))
		return
	}

	tenantID := c.GetString("tenant_id")
	if err := h.service.Delete(c.Request.Context(), id, tenantID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "App component deleted successfully"})
}

// respondError responds with an error
func (h *AppComponentHandler) respondError(c *gin.Context, err error) {
	appErr := errors.FromError(err)
	h.logger.Error("Request failed",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("error", appErr.Message),
	)
	c.JSON(appErr.StatusCode, gin.H{"error": appErr})
}
