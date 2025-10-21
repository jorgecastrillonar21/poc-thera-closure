package services


import (
	"context"
	"fmt"
	"strings"
	"theraclosure/users-service/internal/core/domain"
	"theraclosure/users-service/internal/core/ports"

	"github.com/google/uuid"
)

type userService struct {
	userRepo       ports.UserProfileRepository
	enrollmentRepo ports.EnrollmentRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo ports.UserProfileRepository, enrollmentRepo ports.EnrollmentRepository) ports.UserService {
	return &userService{
		userRepo:       userRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

// CreateProfile creates a new user profile
func (s *userService) CreateProfile(ctx context.Context, profile *domain.UserProfile) error {
	// Validate required fields
	if err := s.ValidateProfile(ctx, profile); err != nil {
		return fmt.Errorf("profile validation failed: %w", err)
	}
	
	// Check if user already exists
	existing, err := s.userRepo.GetByUserID(ctx, profile.UserID)
	if err == nil && existing != nil {
		return fmt.Errorf("profile already exists for user %s", profile.UserID.String())
	}
	
	// Create the profile
	if err := s.userRepo.Create(ctx, profile); err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	
	// Check if profile is complete and update accordingly
	isComplete, _ := s.CheckProfileCompletion(ctx, profile.UserID)
	if isComplete {
		profile.ProfileComplete = true
		s.userRepo.MarkProfileComplete(ctx, profile.UserID)
	}
	
	return nil
}

// GetProfile retrieves a user profile by user ID
func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	profile, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	
	return profile, nil
}

// UpdateProfile updates an existing user profile
func (s *userService) UpdateProfile(ctx context.Context, profile *domain.UserProfile) error {
	// Validate the updated profile
	if err := s.ValidateProfile(ctx, profile); err != nil {
		return fmt.Errorf("profile validation failed: %w", err)
	}
	
	// Update the profile
	if err := s.userRepo.Update(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	
	// Check if profile is now complete
	isComplete, _ := s.CheckProfileCompletion(ctx, profile.UserID)
	if isComplete && !profile.ProfileComplete {
		s.userRepo.MarkProfileComplete(ctx, profile.UserID)
	}
	
	return nil
}

// DeleteProfile deletes a user profile
func (s *userService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	profile, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}
	
	if err := s.userRepo.Delete(ctx, profile.ID); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	
	return nil
}

// ListProfiles lists user profiles with pagination
func (s *userService) ListProfiles(ctx context.Context, limit, offset int) ([]*domain.UserProfile, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	
	profiles, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}
	
	return profiles, nil
}

// ValidateProfile validates a user profile
func (s *userService) ValidateProfile(ctx context.Context, profile *domain.UserProfile) error {
	if strings.TrimSpace(profile.FirstName) == "" {
		return fmt.Errorf("first name is required")
	}
	
	if strings.TrimSpace(profile.LastName) == "" {
		return fmt.Errorf("last name is required")
	}
	
	if strings.TrimSpace(profile.Email) == "" {
		return fmt.Errorf("email is required")
	}
	
	// Validate email format (basic validation)
	if !strings.Contains(profile.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// CheckProfileCompletion checks if a user profile is complete
func (s *userService) CheckProfileCompletion(ctx context.Context, userID uuid.UUID) (bool, error) {
	profile, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get profile: %w", err)
	}
	
	// Define completion criteria
	required := []string{
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.LicenseNumber,
		profile.LicenseState,
		profile.ProfessionalTitle,
	}
	
	for _, field := range required {
		if strings.TrimSpace(field) == "" {
			return false, nil
		}
	}
	
	// Check if practice information is complete
	if profile.PracticeType == "" || profile.PracticeName == "" {
		return false, nil
	}
	
	return true, nil
}

// SearchProfiles searches for profiles by query string
func (s *userService) SearchProfiles(ctx context.Context, query string, limit, offset int) ([]*domain.UserProfile, error) {
	// For now, this is a simple implementation
	// In a real application, you might use full-text search or Elasticsearch
	profiles, err := s.ListProfiles(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	var filtered []*domain.UserProfile
	query = strings.ToLower(strings.TrimSpace(query))
	
	for _, profile := range profiles {
		if s.matchesQuery(profile, query) {
			filtered = append(filtered, profile)
		}
	}
	
	return filtered, nil
}

// GetProfilesByStatus gets profiles by status
func (s *userService) GetProfilesByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.UserProfile, error) {
	// This would typically be implemented at the repository level for better performance
	profiles, err := s.ListProfiles(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	var filtered []*domain.UserProfile
	for _, profile := range profiles {
		if profile.Status == status {
			filtered = append(filtered, profile)
		}
	}
	
	return filtered, nil
}

// matchesQuery checks if a profile matches the search query
func (s *userService) matchesQuery(profile *domain.UserProfile, query string) bool {
	searchFields := []string{
		strings.ToLower(profile.FirstName),
		strings.ToLower(profile.LastName),
		strings.ToLower(profile.Email),
		strings.ToLower(profile.PracticeName),
		strings.ToLower(profile.LicenseNumber),
	}
	
	for _, field := range searchFields {
		if strings.Contains(field, query) {
			return true
		}
	}
	
	return false
}