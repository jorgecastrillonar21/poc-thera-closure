package http

import (
	"net/http"
	"theraclosure/auth-service/internal/adapters/config"
	"theraclosure/auth-service/internal/core/ports"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	authService ports.AuthService
	userService ports.UserService
	config      *config.Config
	router      *gin.Engine
}

// NewServer creates a new HTTP server instance
func NewServer(authService ports.AuthService, userService ports.UserService, config *config.Config) *Server {
	server := &Server{
		authService: authService,
		userService: userService,
		config:      config,
	}

	server.setupRouter()
	return server
}

// setupRouter configures the Gin router with routes and middleware
func (s *Server) setupRouter() {
	// Set Gin mode
	gin.SetMode(s.config.Server.Mode)

	s.router = gin.Default()

	// CORS middleware
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	s.router.GET("/api/v1/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
			auth.GET("/me", s.getCurrentUser)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := s.config.Server.Host + ":" + s.config.Server.Port
	return s.router.Run(addr)
}

// Health check handler
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": s.config.App.Name,
		"version": s.config.App.Version,
		"status":  "healthy",
	})
}

// Placeholder handlers (to be implemented)
func (s *Server) register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Register endpoint not implemented yet"})
}

func (s *Server) login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Login endpoint not implemented yet"})
}

func (s *Server) refreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Refresh token endpoint not implemented yet"})
}

func (s *Server) logout(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Logout endpoint not implemented yet"})
}

func (s *Server) getCurrentUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Get current user endpoint not implemented yet"})
}
