package http


import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"theraclosure/users-service/internal/adapters/config"
	"theraclosure/users-service/internal/core/ports"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config            *config.Config
	userService       ports.UserService
	enrollmentService ports.EnrollmentService
	router            *gin.Engine
	httpServer        *http.Server
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, userService ports.UserService, enrollmentService ports.EnrollmentService) *Server {
	server := &Server{
		config:            cfg,
		userService:       userService,
		enrollmentService: enrollmentService,
	}

	server.setupRouter()
	return server
}

// setupRouter sets up the Gin router with routes and middleware
func (s *Server) setupRouter() {
	// Set Gin mode based on environment
	if s.config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.Default()

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     s.config.CORS.AllowedOrigins,
		AllowMethods:     s.config.CORS.AllowedMethods,
		AllowHeaders:     s.config.CORS.AllowedHeaders,
		AllowCredentials: true,
	}
	s.router.Use(cors.New(corsConfig))

	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// User profile routes
		users := v1.Group("/users")
		{
			users.POST("/profiles", s.createProfile)
			users.GET("/profiles/:userId", s.getProfile)
			users.PUT("/profiles/:userId", s.updateProfile)
			users.DELETE("/profiles/:userId", s.deleteProfile)
			users.GET("/profiles", s.listProfiles)
			users.GET("/profiles/search", s.searchProfiles)
		}

		// Enrollment routes
		enrollments := v1.Group("/enrollments")
		{
			enrollments.POST("/start", s.startEnrollment)
			enrollments.GET("/:userId", s.getEnrollment)
			enrollments.PUT("/:userId", s.updateEnrollment)
			enrollments.POST("/:userId/steps/:step/complete", s.completeStep)
			enrollments.GET("/:userId/progress", s.getProgress)
			enrollments.POST("/:userId/complete", s.completeEnrollment)
			enrollments.PUT("/:userId/plan", s.updatePlan)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Users service starting on %s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give the server 30 seconds to finish handling requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "users-service",
		"timestamp": time.Now().Unix(),
	})
}