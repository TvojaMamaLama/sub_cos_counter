package repository

import (
	"context"
	"fmt"
	"sub-cos-counter/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, req *models.CreatePaymentRequest) (*models.Payment, error) {
	query := `
		INSERT INTO payments (subscription_id, amount, currency, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, subscription_id, amount, currency, paid_at, status, created_at`

	var payment models.Payment
	err := r.db.QueryRow(ctx, query,
		req.SubscriptionID, req.Amount, req.Currency, req.Status,
	).Scan(
		&payment.ID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
		&payment.PaidAt, &payment.Status, &payment.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) GetBySubscriptionID(ctx context.Context, subscriptionID int) ([]*models.Payment, error) {
	query := `
		SELECT id, subscription_id, amount, currency, paid_at, status, created_at
		FROM payments WHERE subscription_id = $1 ORDER BY paid_at DESC`

	rows, err := r.db.Query(ctx, query, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
			&payment.PaidAt, &payment.Status, &payment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, &payment)
	}

	return payments, nil
}

func (r *PaymentRepository) GetMonthlyExpense(ctx context.Context, month time.Time) ([]models.PaymentSummary, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	query := `
		SELECT currency, SUM(amount) as total_amount, COUNT(*) as count
		FROM payments 
		WHERE paid_at >= $1 AND paid_at <= $2 AND status = 'completed'
		GROUP BY currency`

	rows, err := r.db.Query(ctx, query, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly expense: %w", err)
	}
	defer rows.Close()

	var summaries []models.PaymentSummary
	for rows.Next() {
		var summary models.PaymentSummary
		err := rows.Scan(&summary.Currency, &summary.TotalAmount, &summary.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (r *PaymentRepository) GetCategoryAnalytics(ctx context.Context, startDate, endDate time.Time) (map[models.Category][]models.PaymentSummary, error) {
	query := `
		SELECT s.category, p.currency, SUM(p.amount) as total_amount, COUNT(p.*) as count
		FROM payments p
		JOIN subscriptions s ON p.subscription_id = s.id
		WHERE p.paid_at >= $1 AND p.paid_at <= $2 AND p.status = 'completed'
		GROUP BY s.category, p.currency
		ORDER BY s.category, p.currency`

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get category analytics: %w", err)
	}
	defer rows.Close()

	analytics := make(map[models.Category][]models.PaymentSummary)
	for rows.Next() {
		var category models.Category
		var summary models.PaymentSummary
		err := rows.Scan(&category, &summary.Currency, &summary.TotalAmount, &summary.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category analytics: %w", err)
		}
		analytics[category] = append(analytics[category], summary)
	}

	return analytics, nil
}

func (r *PaymentRepository) GetAllPayments(ctx context.Context, limit int) ([]*models.Payment, error) {
	query := `
		SELECT id, subscription_id, amount, currency, paid_at, status, created_at
		FROM payments ORDER BY paid_at DESC LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get all payments: %w", err)
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
			&payment.PaidAt, &payment.Status, &payment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, &payment)
	}

	return payments, nil
}
