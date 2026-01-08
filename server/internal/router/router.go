package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-system-config-service/internal/handler"
)

// SetupRouter sets up the Gin router with all routes
func SetupRouter(
	appComponentHandler *handler.AppComponentHandler,
	countryHandler *handler.CountryHandler,
	log *logger.Logger,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "system-config-service",
		})
	})
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"service": "system-config-service",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1/system-config")
	{
		// App Components
		appComponents := v1.Group("/app-components")
		{
			appComponents.GET("", appComponentHandler.List)
			appComponents.GET("/:id", appComponentHandler.GetByID)
			appComponents.POST("", appComponentHandler.Create)
			appComponents.PUT("/:id", appComponentHandler.Update)
			appComponents.DELETE("/:id", appComponentHandler.Delete)
		}

		// Countries
		countries := v1.Group("/countries")
		{
			countries.GET("", countryHandler.List)
			countries.GET("/:code", countryHandler.GetByCode)
			countries.POST("", countryHandler.Create)
			countries.PUT("/:code", countryHandler.Update)
			countries.DELETE("/:code", countryHandler.Delete)
		}

		// Placeholder routes for other entities
		// These would be implemented similarly to the above

		// SaaS Modules
		modules := v1.Group("/modules")
		{
			modules.GET("", placeholderHandler)
			modules.GET("/:id", placeholderHandler)
			modules.POST("", placeholderHandler)
			modules.PUT("/:id", placeholderHandler)
			modules.DELETE("/:id", placeholderHandler)
		}

		// Service Packages
		packages := v1.Group("/packages")
		{
			packages.GET("", placeholderHandler)
			packages.GET("/:id", placeholderHandler)
			packages.POST("", placeholderHandler)
			packages.PUT("/:id", placeholderHandler)
			packages.DELETE("/:id", placeholderHandler)
		}

		// Admin Menus
		menus := v1.Group("/menus")
		{
			menus.GET("", placeholderHandler)
			menus.GET("/tree", placeholderHandler)
			menus.GET("/by-module/:module_code", placeholderHandler)
			menus.GET("/:id", placeholderHandler)
			menus.POST("", placeholderHandler)
			menus.PUT("/:id", placeholderHandler)
			menus.DELETE("/:id", placeholderHandler)
		}

		// Permissions
		permissions := v1.Group("/permissions")
		{
			permissions.GET("", placeholderHandler)
			permissions.GET("/:id", placeholderHandler)
			permissions.GET("/by-module/:module_code", placeholderHandler)
			permissions.GET("/by-resource/:resource", placeholderHandler)
			permissions.POST("", placeholderHandler)
			permissions.PUT("/:id", placeholderHandler)
			permissions.DELETE("/:id", placeholderHandler)
			permissions.POST("/batch", placeholderHandler)
		}

		// Roles
		roles := v1.Group("/roles")
		{
			roles.GET("", placeholderHandler)
			roles.GET("/:id", placeholderHandler)
			roles.POST("", placeholderHandler)
			roles.PUT("/:id", placeholderHandler)
			roles.DELETE("/:id", placeholderHandler)
			roles.GET("/:id/permissions", placeholderHandler)
			roles.PUT("/:id/permissions", placeholderHandler)
			roles.POST("/:id/clone", placeholderHandler)
		}

		// Ethnicities
		ethnicities := v1.Group("/ethnicities")
		{
			ethnicities.GET("", placeholderHandler)
			ethnicities.GET("/:id", placeholderHandler)
			ethnicities.GET("/by-country/:country_code", placeholderHandler)
			ethnicities.POST("", placeholderHandler)
			ethnicities.PUT("/:id", placeholderHandler)
			ethnicities.DELETE("/:id", placeholderHandler)
		}

		// Locations (Hierarchical)
		locations := v1.Group("/locations")
		{
			locations.GET("/countries/:country_code/provinces", placeholderHandler)
			locations.GET("/provinces/:province_code", placeholderHandler)
			locations.GET("/provinces/:province_code/districts", placeholderHandler)
			locations.GET("/districts/:district_code", placeholderHandler)
			locations.GET("/districts/:district_code/wards", placeholderHandler)
			locations.GET("/wards/:ward_code", placeholderHandler)
			locations.GET("/search", placeholderHandler)
			locations.POST("/provinces", placeholderHandler)
			locations.POST("/districts", placeholderHandler)
			locations.POST("/wards", placeholderHandler)
			locations.PUT("/provinces/:code", placeholderHandler)
			locations.PUT("/districts/:code", placeholderHandler)
			locations.PUT("/wards/:code", placeholderHandler)
			locations.DELETE("/provinces/:code", placeholderHandler)
			locations.DELETE("/districts/:code", placeholderHandler)
			locations.DELETE("/wards/:code", placeholderHandler)
		}

		// Currencies
		currencies := v1.Group("/currencies")
		{
			currencies.GET("", placeholderHandler)
			currencies.GET("/:code", placeholderHandler)
			currencies.POST("", placeholderHandler)
			currencies.PUT("/:code", placeholderHandler)
			currencies.DELETE("/:code", placeholderHandler)
		}
	}

	return router
}

// placeholderHandler is a temporary handler for routes that are not yet implemented
func placeholderHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "This endpoint is not yet implemented",
		"path":    c.Request.URL.Path,
	})
}
