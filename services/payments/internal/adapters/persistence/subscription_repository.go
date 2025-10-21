package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/ports"
)

type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new subscription repository instance
func NewSubscriptionRepository(database *Database) ports.SubscriptionRepository {
	return &subscriptionRepository{
		db: database.GetDB(),
	}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	if err := r.db.WithContext(ctx).Create(subscription).Error; err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	var subscription domain.Subscription
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Payments").
		First(&subscription, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return &subscription, nil
}

func (r *subscriptionRepository) GetByStripeID(ctx context.Context, stripeID string) (*domain.Subscription, error) {
	var subscription domain.Subscription
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Payments").
		First(&subscription, "stripe_id = ?", stripeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription by stripe ID: %w", err)
	}
	return &subscription, nil
}

func (r *subscriptionRepository) GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Subscription, error) {
	var subscriptions []*domain.Subscription
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Payments").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by customer ID: %w", err)
	}
	return subscriptions, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	if err := r.db.WithContext(ctx).Save(subscription).Error; err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Subscription{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) List(ctx context.Context, offset, limit int) ([]*domain.Subscription, int64, error) {
	var subscriptions []*domain.Subscription
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Subscription{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	// Get subscriptions with pagination
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Payments").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&subscriptions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subscriptions, total, nil
}