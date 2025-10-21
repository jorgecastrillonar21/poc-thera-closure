package domain


import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Country represents a country entity
type Country struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Code      string    `json:"code" gorm:"uniqueIndex;not null;size:3"` // ISO 3166-1 alpha-3 code (e.g., USA, CAN, MEX)
	Code2     string    `json:"code2" gorm:"uniqueIndex;not null;size:2"` // ISO 3166-1 alpha-2 code (e.g., US, CA, MX)
	Region    string    `json:"region" gorm:"size:50"`                    // North America, Europe, etc.
	Currency  string    `json:"currency" gorm:"size:10"`                  // USD, CAD, EUR, etc.
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	States []State `json:"states,omitempty" gorm:"foreignKey:CountryID"`
}

// State represents a state/province/region entity
type State struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CountryID string    `json:"country_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Code      string    `json:"code" gorm:"size:10"` // State/Province code (e.g., CA, NY, ON, etc.)
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Country Country `json:"country,omitempty" gorm:"foreignKey:CountryID"`
	Cities  []City  `json:"cities,omitempty" gorm:"foreignKey:StateID"`
}

// City represents a city entity
type City struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StateID   string    `json:"state_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	ZipCode   string    `json:"zip_code" gorm:"size:20"`  // Primary zip/postal code
	Latitude  float64   `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude float64   `json:"longitude" gorm:"type:decimal(11,8)"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	State State `json:"state,omitempty" gorm:"foreignKey:StateID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Country) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (s *State) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (c *City) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table names
func (Country) TableName() string {
	return "countries"
}

func (State) TableName() string {
	return "states"
}

func (City) TableName() string {
	return "cities"
}

// Validation methods
func (c *Country) IsValid() bool {
	return c.Name != "" && c.Code != "" && c.Code2 != ""
}

func (s *State) IsValid() bool {
	return s.Name != "" && s.CountryID != ""
}

func (ci *City) IsValid() bool {
	return ci.Name != "" && ci.StateID != ""
}