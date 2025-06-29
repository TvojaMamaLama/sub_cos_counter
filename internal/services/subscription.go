package services

import (
	"context"
	"fmt"
	"sub-cos-counter/internal/models"
	"sub-cos-counter/internal/repository"
)

type SubscriptionService struct {
	subscriptionRepo *repository.SubscriptionRepository
	paymentRepo      *repository.PaymentRepository
}

func NewSubscriptionService(subscriptionRepo *repository.SubscriptionRepository, paymentRepo *repository.PaymentRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		paymentRepo:      paymentRepo,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("subscription name is required")
	}
	if req.Cost <= 0 {
		return nil, fmt.Errorf("subscription cost must be positive")
	}
	if req.PeriodDays <= 0 {
		return nil, fmt.Errorf("subscription period must be positive")
	}

	return s.subscriptionRepo.Create(ctx, req)
}

func (s *SubscriptionService) GetAllActiveSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	return s.subscriptionRepo.GetAllActive(ctx)
}

func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error) {
	return s.subscriptionRepo.GetByID(ctx, id)
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return s.subscriptionRepo.Delete(ctx, id)
}

func (s *SubscriptionService) MarkAsPaid(ctx context.Context, subscriptionID int) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Create payment record
	paymentReq := &models.CreatePaymentRequest{
		SubscriptionID: subscriptionID,
		Amount:         subscription.Cost,
		Currency:       subscription.Currency,
		Status:         models.PaymentStatusCompleted,
	}

	_, err = s.paymentRepo.Create(ctx, paymentReq)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// Update next payment date
	subscription.UpdateNextPayment()
	err = s.subscriptionRepo.Update(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (s *SubscriptionService) GetDuePayments(ctx context.Context) ([]*models.Subscription, error) {
	return s.subscriptionRepo.GetDuePayments(ctx)
}

func (s *SubscriptionService) GetSubscriptionsByCategory(ctx context.Context, category models.Category) ([]*models.Subscription, error) {
	return s.subscriptionRepo.GetByCategory(ctx, category)
}