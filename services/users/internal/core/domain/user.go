package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserProfile represents extended user profile information
type UserProfile struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`

	// Personal Information
	FirstName   string     `json:"first_name" gorm:"size:100;not null"`
	LastName    string     `json:"last_name" gorm:"size:100;not null"`
	Email       string     `json:"email" gorm:"size:255;not null;uniqueIndex"`
	Phone       string     `json:"phone" gorm:"size:20"`
	DateOfBirth *time.Time `json:"date_of_birth"`

	// Address Information
	Address string `json:"address" gorm:"size:255"`
	City    string `json:"city" gorm:"size:100"`
	State   string `json:"state" gorm:"size:50"`
	ZipCode string `json:"zip_code" gorm:"size:20"`
	Country string `json:"country" gorm:"size:100;default:'US'"`

	// Professional Information
	LicenseNumber     string     `json:"license_number" gorm:"size:100"`
	LicenseState      string     `json:"license_state" gorm:"size:50"`
	LicenseExpiration *time.Time `json:"license_expiration"`
	ProfessionalTitle string     `json:"professional_title" gorm:"size:100"`
	Specializations   []string   `json:"specializations" gorm:"type:text[]"`
	YearsOfExperience int        `json:"years_of_experience"`

	// Practice Information
	PracticeName    string `json:"practice_name" gorm:"size:255"`
	PracticeType    string `json:"practice_type" gorm:"size:100"` // individual, group, clinic, hospital
	PracticeAddress string `json:"practice_address" gorm:"size:255"`
	PracticeCity    string `json:"practice_city" gorm:"size:100"`
	PracticeState   string `json:"practice_state" gorm:"size:50"`
	PracticeZipCode string `json:"practice_zip_code" gorm:"size:20"`
	PracticePhone   string `json:"practice_phone" gorm:"size:20"`

	// Additional Information
	EmergencyContactName  string `json:"emergency_contact_name" gorm:"size:255"`
	EmergencyContactPhone string `json:"emergency_contact_phone" gorm:"size:20"`
	EmergencyContactEmail string `json:"emergency_contact_email" gorm:"size:255"`

	// Status and Metadata
	ProfileComplete bool   `json:"profile_complete" gorm:"default:false"`
	Status          string `json:"status" gorm:"size:50;default:'active'"` // active, inactive, suspended

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// EnrollmentData represents user enrollment progress and data
type EnrollmentData struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`

	// Enrollment Steps Progress
	PersonalInfoComplete     bool `json:"personal_info_complete" gorm:"default:false"`
	LicensureDetailsComplete bool `json:"licensure_details_complete" gorm:"default:false"`
	PracticeInfoComplete     bool `json:"practice_info_complete" gorm:"default:false"`
	AdminSetupComplete       bool `json:"admin_setup_complete" gorm:"default:false"`
	ScheduleConfigComplete   bool `json:"schedule_config_complete" gorm:"default:false"`

	// Enrollment Status
	EnrollmentStatus string     `json:"enrollment_status" gorm:"size:50;default:'in_progress'"` // in_progress, completed, paused
	CurrentStep      int        `json:"current_step" gorm:"default:1"`
	TotalSteps       int        `json:"total_steps" gorm:"default:5"`
	CompletionDate   *time.Time `json:"completion_date"`

	// Plan Information
	SelectedPlan  string `json:"selected_plan" gorm:"size:50"`                    // essential, professional, enterprise
	PaymentStatus string `json:"payment_status" gorm:"size:50;default:'pending'"` // pending, completed, failed

	// Additional Enrollment Data
	PreferredContactMethod string `json:"preferred_contact_method" gorm:"size:50;default:'email'"` // email, phone, text
	ReferralSource         string `json:"referral_source" gorm:"size:255"`
	SpecialRequests        string `json:"special_requests" gorm:"type:text"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}

// TableName specifies the table name for EnrollmentData
func (EnrollmentData) TableName() string {
	return "enrollment_data"
}
