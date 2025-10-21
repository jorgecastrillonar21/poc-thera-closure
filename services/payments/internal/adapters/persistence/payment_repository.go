package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/ports"
)

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository instance
func NewPaymentRepository(database *Database) ports.PaymentRepository {
	return &paymentRepository{
		db: database.GetDB(),
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	if err := r.db.WithContext(ctx).Create(payment).Error; err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Subscription").
		First(&payment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return &payment, nil
}

func (r *paymentRepository) GetByStripeID(ctx context.Context, stripeID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Subscription").
		First(&payment, "stripe_id = ?", stripeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment by stripe ID: %w", err)
	}
	return &payment, nil
}

func (r *paymentRepository) GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Subscription").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get payments by customer ID: %w", err)
	}
	return payments, nil
}

func (r *paymentRepository) GetBySubscriptionID(ctx context.Context, subscriptionID string) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Subscription").
		Where("subscription_id = ?", subscriptionID).
		Order("created_at DESC").
		Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get payments by subscription ID: %w", err)
	}
	return payments, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	if err := r.db.WithContext(ctx).Save(payment).Error; err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Payment{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) List(ctx context.Context, offset, limit int) ([]*domain.Payment, int64, error) {
	var payments []*domain.Payment
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Get payments with pagination
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Subscription").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&payments).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}

	return payments, total, nil
}
