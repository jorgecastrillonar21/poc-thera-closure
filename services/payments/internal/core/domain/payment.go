package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusUnpaid    SubscriptionStatus = "unpaid"
	SubscriptionStatusPaused    SubscriptionStatus = "paused"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCanceled  PaymentStatus = "canceled"
)

// Customer represents a customer entity
type Customer struct {
	ID             string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         string    `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"` // Reference to users service
	StripeID       string    `json:"stripe_id" gorm:"uniqueIndex;size:255"`
	Email          string    `json:"email" gorm:"not null;size:255"`
	Name           string    `json:"name" gorm:"not null;size:255"`
	DefaultPaymentMethodID string `json:"default_payment_method_id" gorm:"size:255"`
	Active         bool      `json:"active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Subscriptions []Subscription `json:"subscriptions,omitempty" gorm:"foreignKey:CustomerID"`
	Payments      []Payment      `json:"payments,omitempty" gorm:"foreignKey:CustomerID"`
}

// Subscription represents a subscription entity
type Subscription struct {
	ID                 string             `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID         string             `json:"customer_id" gorm:"type:uuid;not null"`
	StripeID           string             `json:"stripe_id" gorm:"uniqueIndex;size:255"`
	PriceID            string             `json:"price_id" gorm:"not null;size:255"` // Stripe Price ID
	Status             SubscriptionStatus `json:"status" gorm:"not null"`
	CurrentPeriodStart time.Time          `json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end"`
	TrialEnd           *time.Time         `json:"trial_end"`
	CancelAt           *time.Time         `json:"cancel_at"`
	CanceledAt         *time.Time         `json:"canceled_at"`
	Amount             int64              `json:"amount"` // Amount in cents
	Currency           string             `json:"currency" gorm:"size:3;default:'usd'"`
	Active             bool               `json:"active" gorm:"default:true"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	DeletedAt          gorm.DeletedAt     `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:SubscriptionID"`
}

// Payment represents a payment entity
type Payment struct {
	ID             string        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID     string        `json:"customer_id" gorm:"type:uuid;not null"`
	SubscriptionID *string       `json:"subscription_id" gorm:"type:uuid"` // Optional for one-time payments
	StripeID       string        `json:"stripe_id" gorm:"uniqueIndex;size:255"`
	PaymentIntentID string       `json:"payment_intent_id" gorm:"size:255"`
	Status         PaymentStatus `json:"status" gorm:"not null"`
	Amount         int64         `json:"amount"` // Amount in cents
	Currency       string        `json:"currency" gorm:"size:3;default:'usd'"`
	Description    string        `json:"description" gorm:"size:500"`
	Metadata       string        `json:"metadata" gorm:"type:jsonb"` // JSON metadata
	ProcessedAt    *time.Time    `json:"processed_at"`
	RefundedAt     *time.Time    `json:"refunded_at"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Customer     Customer     `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table names
func (Customer) TableName() string {
	return "customers"
}

func (Subscription) TableName() string {
	return "subscriptions"
}

func (Payment) TableName() string {
	return "payments"
}

// Validation methods
func (c *Customer) IsValid() bool {
	return c.UserID != "" && c.Email != "" && c.Name != ""
}

func (s *Subscription) IsValid() bool {
	return s.CustomerID != "" && s.PriceID != "" && s.Status != ""
}

func (p *Payment) IsValid() bool {
	return p.CustomerID != "" && p.Amount > 0 && p.Currency != ""
}