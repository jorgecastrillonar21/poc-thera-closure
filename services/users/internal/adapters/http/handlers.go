package http

import (
	"net/http"
	"strconv"
	"theraclosure/users-service/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User Profile Handlers

// createProfile handles creating a new user profile
func (s *Server) createProfile(c *gin.Context) {
	var profile domain.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Auto-generate UserID if not provided
	if profile.UserID == uuid.Nil {
		profile.UserID = uuid.New()
	}

	if err := s.userService.CreateProfile(c.Request.Context(), &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile created successfully",
		"profile": profile,
	})
}

// getProfile handles getting a user profile by user ID
func (s *Server) getProfile(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	profile, err := s.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Profile not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
	})
}

// updateProfile handles updating a user profile
func (s *Server) updateProfile(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get the existing profile first
	existingProfile, err := s.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Profile not found",
			"details": err.Error(),
		})
		return
	}

	// Bind the JSON to the existing profile
	if err := c.ShouldBindJSON(existingProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Ensure the user ID and ID remain unchanged
	existingProfile.UserID = userID

	if err := s.userService.UpdateProfile(c.Request.Context(), existingProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": existingProfile,
	})
}

// deleteProfile handles deleting a user profile
func (s *Server) deleteProfile(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	if err := s.userService.DeleteProfile(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete profile",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// listProfiles handles listing user profiles with pagination
func (s *Server) listProfiles(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid offset parameter",
		})
		return
	}

	profiles, err := s.userService.ListProfiles(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list profiles",
			"details": err.Error(),
		})
		return
	}

	// Get total count (simplified - in production this would be optimized)
	totalProfiles, _ := s.userService.ListProfiles(c.Request.Context(), 1000000, 0) // Large number to get all
	totalCount := len(totalProfiles)

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"limit":    limit,
		"offset":   offset,
		"count":    len(profiles),
		"total":    totalCount,
	})
}

// searchProfiles handles searching user profiles
func (s *Server) searchProfiles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid offset parameter",
		})
		return
	}

	profiles, err := s.userService.SearchProfiles(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search profiles",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"query":    query,
		"limit":    limit,
		"offset":   offset,
		"count":    len(profiles),
	})
}

// Enrollment Handlers

// startEnrollment handles starting a new enrollment
func (s *Server) startEnrollment(c *gin.Context) {
	var request struct {
		UserID       string `json:"user_id" binding:"required"`
		SelectedPlan string `json:"selected_plan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	if err := s.enrollmentService.StartEnrollment(c.Request.Context(), userID, request.SelectedPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to start enrollment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Enrollment started successfully",
		"user_id": request.UserID,
		"plan":    request.SelectedPlan,
	})
}

// getEnrollment handles getting enrollment data
func (s *Server) getEnrollment(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	enrollment, err := s.enrollmentService.GetEnrollment(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Enrollment not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollment": enrollment,
	})
}

// updateEnrollment handles updating enrollment data
func (s *Server) updateEnrollment(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var enrollment domain.EnrollmentData
	if err := c.ShouldBindJSON(&enrollment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Ensure the user ID matches
	enrollment.UserID = userID

	if err := s.enrollmentService.UpdateEnrollment(c.Request.Context(), &enrollment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update enrollment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Enrollment updated successfully",
		"enrollment": enrollment,
	})
}

// completeStep handles completing an enrollment step
func (s *Server) completeStep(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	stepStr := c.Param("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid step format",
		})
		return
	}

	if err := s.enrollmentService.CompleteStep(c.Request.Context(), userID, step); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to complete step",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Step completed successfully",
		"user_id": userID,
		"step":    step,
	})
}

// getProgress handles getting enrollment progress
func (s *Server) getProgress(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	progress, err := s.enrollmentService.GetProgress(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get progress",
			"details": err.Error(),
		})
		return
	}

	currentStep, err := s.enrollmentService.GetCurrentStep(c.Request.Context(), userID)
	if err != nil {
		currentStep = 1 // Default to step 1 if error
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"progress":     progress,
		"current_step": currentStep,
		"total_steps":  5,
	})
}

// completeEnrollment handles completing the entire enrollment
func (s *Server) completeEnrollment(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	if err := s.enrollmentService.CompleteEnrollment(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to complete enrollment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Enrollment completed successfully",
		"user_id": userID,
	})
}

// updatePlan handles updating the selected plan
func (s *Server) updatePlan(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var request struct {
		Plan string `json:"plan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := s.enrollmentService.UpdateSelectedPlan(c.Request.Context(), userID, request.Plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plan updated successfully",
		"user_id": userID,
		"plan":    request.Plan,
	})
}