package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/ports"
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository instance
func NewCustomerRepository(database *Database) ports.CustomerRepository {
	return &customerRepository{
		db: database.GetDB(),
	}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	if err := r.db.WithContext(ctx).Create(customer).Error; err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.WithContext(ctx).
		Preload("Subscriptions").
		Preload("Payments").
		First(&customer, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

func (r *customerRepository) GetByUserID(ctx context.Context, userID string) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.WithContext(ctx).
		Preload("Subscriptions").
		Preload("Payments").
		First(&customer, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer by user ID: %w", err)
	}
	return &customer, nil
}

func (r *customerRepository) GetByStripeID(ctx context.Context, stripeID string) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.WithContext(ctx).
		Preload("Subscriptions").
		Preload("Payments").
		First(&customer, "stripe_id = ?", stripeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer by stripe ID: %w", err)
	}
	return &customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	if err := r.db.WithContext(ctx).Save(customer).Error; err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Customer{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}

func (r *customerRepository) List(ctx context.Context, offset, limit int) ([]*domain.Customer, int64, error) {
	var customers []*domain.Customer
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Get customers with pagination
	if err := r.db.WithContext(ctx).
		Preload("Subscriptions").
		Preload("Payments").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&customers).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}

	return customers, total, nil
}