package repository

import (
	"context"
	"fmt"
	"sub-cos-counter/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	query := `
		INSERT INTO subscriptions (name, cost, currency, period_days, next_payment, category, auto_renewal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, cost, currency, period_days, next_payment, category, auto_renewal, active, created_at, updated_at`

	var sub models.Subscription
	err := r.db.QueryRow(ctx, query,
		req.Name, req.Cost, req.Currency, req.PeriodDays, req.NextPayment, req.Category, req.AutoRenewal,
	).Scan(
		&sub.ID, &sub.Name, &sub.Cost, &sub.Currency, &sub.PeriodDays,
		&sub.NextPayment, &sub.Category, &sub.AutoRenewal, &sub.Active, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `
		SELECT id, name, cost, currency, period_days, next_payment, category, auto_renewal, active, created_at, updated_at
		FROM subscriptions WHERE id = $1`

	var sub models.Subscription
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.Name, &sub.Cost, &sub.Currency, &sub.PeriodDays,
		&sub.NextPayment, &sub.Category, &sub.AutoRenewal, &sub.Active, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetAllActive(ctx context.Context) ([]*models.Subscription, error) {
	query := `
		SELECT id, name, cost, currency, period_days, next_payment, category, auto_renewal, active, created_at, updated_at
		FROM subscriptions WHERE active = true ORDER BY next_payment ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Cost, &sub.Currency, &sub.PeriodDays,
			&sub.NextPayment, &sub.Category, &sub.AutoRenewal, &sub.Active, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET name = $2, cost = $3, currency = $4, period_days = $5, next_payment = $6, 
		    category = $7, auto_renewal = $8, active = $9, updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query,
		sub.ID, sub.Name, sub.Cost, sub.Currency, sub.PeriodDays,
		sub.NextPayment, sub.Category, sub.AutoRenewal, sub.Active,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE subscriptions SET active = false, updated_at = NOW() WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) GetByCategory(ctx context.Context, category models.Category) ([]*models.Subscription, error) {
	query := `
		SELECT id, name, cost, currency, period_days, next_payment, category, auto_renewal, active, created_at, updated_at
		FROM subscriptions WHERE category = $1 AND active = true ORDER BY cost DESC`

	rows, err := r.db.Query(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by category: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Cost, &sub.Currency, &sub.PeriodDays,
			&sub.NextPayment, &sub.Category, &sub.AutoRenewal, &sub.Active, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) GetDuePayments(ctx context.Context) ([]*models.Subscription, error) {
	query := `
		SELECT id, name, cost, currency, period_days, next_payment, category, auto_renewal, active, created_at, updated_at
		FROM subscriptions 
		WHERE active = true AND next_payment <= $1 
		ORDER BY next_payment ASC`

	rows, err := r.db.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get due payments: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Cost, &sub.Currency, &sub.PeriodDays,
			&sub.NextPayment, &sub.Category, &sub.AutoRenewal, &sub.Active, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}
