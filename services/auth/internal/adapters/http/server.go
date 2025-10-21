package http

import (
	"net/http"
	"theraclosure/auth-service/internal/adapters/config"
	"theraclosure/auth-service/internal/core/domain"
	"theraclosure/auth-service/internal/core/ports"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Server represents the HTTP server
type Server struct {
	authService ports.AuthService
	userService ports.UserService
	config      *config.Config
	router      *gin.Engine
}

// HTTP Request/Response Types
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	SessionID string `json:"sessionId"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresAt    time.Time    `json:"expiresAt"`
}

type TokenResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
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

// HTTP Handlers
func (s *Server) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Missing required fields: email, password, firstName, lastName",
		})
		return
	}

	// Convert to domain request
	domainReq := &domain.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Call auth service
	result, err := s.authService.Register(c.Request.Context(), domainReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convert to response
	response := AuthResponse{
		User: UserResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			FirstName: result.User.FirstName,
			LastName:  result.User.LastName,
			Role:      string(result.User.Role),
			IsActive:  result.User.IsActive,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenDuration),
	}

	c.JSON(http.StatusCreated, response)
}

func (s *Server) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Missing required fields: email, password",
		})
		return
	}

	// Convert to domain request
	domainReq := &domain.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call auth service
	result, err := s.authService.Login(c.Request.Context(), domainReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convert to response
	response := AuthResponse{
		User: UserResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			FirstName: result.User.FirstName,
			LastName:  result.User.LastName,
			Role:      string(result.User.Role),
			IsActive:  result.User.IsActive,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenDuration),
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Missing required field: refreshToken",
		})
		return
	}

	// Convert to domain request
	domainReq := &domain.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}

	// Call auth service
	result, err := s.authService.RefreshToken(c.Request.Context(), domainReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Return token pair
	response := TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenDuration),
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) logout(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Missing authorization header",
		})
		return
	}

	// Remove "Bearer " prefix
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid authorization header format",
		})
		return
	}

	// Validate token and get claims to extract session ID
	_, err := s.authService.ValidateToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// For now, we'll need to get the session ID from the token claims
	// This is a workaround - ideally we'd extract it from JWT claims directly
	// But we need to get the session ID from the token, so let's accept it in request body as backup
	var req LogoutRequest
	var sessionID uuid.UUID
	
	if err := c.ShouldBindJSON(&req); err == nil && req.SessionID != "" {
		sessionID, err = uuid.Parse(req.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid session ID format",
			})
			return
		}
	} else {
		// If no session ID provided, we'll invalidate all sessions for this user
		// This is a simple approach for now
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Session ID required in request body",
		})
		return
	}

	// Call auth service
	if err := s.authService.Logout(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (s *Server) getCurrentUser(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Missing authorization header",
		})
		return
	}

	// Remove "Bearer " prefix
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid authorization header format",
		})
		return
	}

	// Validate token and get user
	user, err := s.authService.ValidateToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convert to response
	response := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}
