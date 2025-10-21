package http


import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"theraclosure/geolocation-service/internal/adapters/config"
	"theraclosure/geolocation-service/internal/core/ports"
)

type Server struct {
	service ports.GeolocationService
	config  *config.Config
	router  *gin.Engine
}

func NewServer(service ports.GeolocationService, cfg *config.Config) *Server {
	server := &Server{
		service: service,
		config:  cfg,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	// Set Gin mode based on config
	if s.config.App.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.New()
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     s.config.CORS.AllowedOrigins,
		AllowMethods:     s.config.CORS.AllowedMethods,
		AllowHeaders:     s.config.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	s.router.Use(cors.New(corsConfig))

	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("/health", s.healthCheck)
		
		// Country routes
		countries := v1.Group("/countries")
		{
			countries.POST("", s.createCountry)
			countries.GET("", s.listCountries)
			countries.GET("/search", s.searchCountries)
			countries.GET("/:id", s.getCountry)
			countries.PUT("/:id", s.updateCountry)
			countries.DELETE("/:id", s.deleteCountry)
			countries.GET("/:id/states", s.listStatesByCountry)
			countries.GET("/:id/hierarchy", s.getCountryHierarchy)
		}

		// State routes
		states := v1.Group("/states")
		{
			states.POST("", s.createState)
			states.GET("/search", s.searchStates)
			states.GET("/:id", s.getState)
			states.PUT("/:id", s.updateState)
			states.DELETE("/:id", s.deleteState)
			states.GET("/:id/cities", s.listCitiesByState)
			states.GET("/:id/hierarchy", s.getStateHierarchy)
		}

		// City routes
		cities := v1.Group("/cities")
		{
			cities.POST("", s.createCity)
			cities.GET("/search", s.searchCities)
			cities.GET("/:id", s.getCity)
			cities.PUT("/:id", s.updateCity)
			cities.DELETE("/:id", s.deleteCity)
		}

		// Hierarchy routes
		hierarchy := v1.Group("/hierarchy")
		{
			hierarchy.GET("", s.getCompleteHierarchy)
		}

		// Bulk operations
		bulk := v1.Group("/bulk")
		{
			bulk.POST("/countries", s.bulkCreateCountries)
			bulk.POST("/states", s.bulkCreateStates)
			bulk.POST("/cities", s.bulkCreateCities)
		}
	}
}

func (s *Server) Start() error {
	return s.router.Run(s.config.GetServerAddress())
}

// Health check handler
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":   s.config.App.Name,
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// Helper function to get pagination parameters
func (s *Server) getPaginationParams(c *gin.Context) (limit, offset int) {
	limit = 20 // default
	offset = 0 // default

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

// Helper function to handle errors
func (s *Server) handleError(c *gin.Context, err error, defaultStatus int) {
	if err.Error() == "not found" || err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	c.JSON(defaultStatus, gin.H{"error": err.Error()})
}