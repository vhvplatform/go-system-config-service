package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/errors"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-system-config-service/internal/domain"
	"github.com/vhvplatform/go-system-config-service/internal/service"
	"go.uber.org/zap"
)

// CountryHandler handles HTTP requests for countries
type CountryHandler struct {
	service *service.CountryService
	logger  *logger.Logger
}

// NewCountryHandler creates a new country handler
func NewCountryHandler(service *service.CountryService, log *logger.Logger) *CountryHandler {
	return &CountryHandler{
		service: service,
		logger:  log,
	}
}

// Create handles creating a new country
func (h *CountryHandler) Create(c *gin.Context) {
	var country domain.Country
	if err := c.ShouldBindJSON(&country); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	if err := h.service.Create(c.Request.Context(), &country); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": country})
}

// GetByCode handles getting a country by code
func (h *CountryHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		h.respondError(c, errors.BadRequest("Code is required"))
		return
	}

	country, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": country})
}

// List handles listing countries
func (h *CountryHandler) List(c *gin.Context) {
	var req domain.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PerPage = 30
	}

	countries, total, err := h.service.List(c.Request.Context(), req.Page, req.PerPage)
	if err != nil {
		h.respondError(c, err)
		return
	}

	totalPages := int(total) / req.PerPage
	if int(total)%req.PerPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": countries,
		"pagination": domain.PaginationResponse{
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: totalPages,
			TotalItems: total,
		},
	})
}

// Update handles updating a country
func (h *CountryHandler) Update(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		h.respondError(c, errors.BadRequest("Code is required"))
		return
	}

	var country domain.Country
	if err := c.ShouldBindJSON(&country); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	country.Code = code

	if err := h.service.Update(c.Request.Context(), &country); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": country})
}

// Delete handles deleting a country
func (h *CountryHandler) Delete(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		h.respondError(c, errors.BadRequest("Code is required"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), code); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Country deleted successfully"})
}

// respondError responds with an error
func (h *CountryHandler) respondError(c *gin.Context, err error) {
	appErr := errors.FromError(err)
	h.logger.Error("Request failed",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("error", appErr.Message),
	)
	c.JSON(appErr.StatusCode, gin.H{"error": appErr})
}
