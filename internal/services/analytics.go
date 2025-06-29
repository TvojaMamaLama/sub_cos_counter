package services

import (
	"context"
	"sub-cos-counter/internal/models"
	"sub-cos-counter/internal/repository"
	"time"
)

type AnalyticsService struct {
	paymentRepo      *repository.PaymentRepository
	subscriptionRepo *repository.SubscriptionRepository
}

func NewAnalyticsService(paymentRepo *repository.PaymentRepository, subscriptionRepo *repository.SubscriptionRepository) *AnalyticsService {
	return &AnalyticsService{
		paymentRepo:      paymentRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *AnalyticsService) GetMonthlyExpense(ctx context.Context, month time.Time) ([]models.PaymentSummary, error) {
	return s.paymentRepo.GetMonthlyExpense(ctx, month)
}

func (s *AnalyticsService) GetCurrentMonthExpense(ctx context.Context) ([]models.PaymentSummary, error) {
	now := time.Now()
	return s.GetMonthlyExpense(ctx, now)
}

func (s *AnalyticsService) GetCategoryAnalytics(ctx context.Context, startDate, endDate time.Time) (map[models.Category][]models.PaymentSummary, error) {
	return s.paymentRepo.GetCategoryAnalytics(ctx, startDate, endDate)
}

func (s *AnalyticsService) GetCurrentMonthCategoryAnalytics(ctx context.Context) (map[models.Category][]models.PaymentSummary, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	
	return s.GetCategoryAnalytics(ctx, startOfMonth, endOfMonth)
}

func (s *AnalyticsService) GetPaymentHistory(ctx context.Context, limit int) ([]*models.Payment, error) {
	return s.paymentRepo.GetAllPayments(ctx, limit)
}

func (s *AnalyticsService) GetUpcomingPayments(ctx context.Context, days int) ([]*models.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().AddDate(0, 0, days)
	var upcoming []*models.Subscription

	for _, sub := range subscriptions {
		if sub.NextPayment.Before(cutoff) {
			upcoming = append(upcoming, sub)
		}
	}

	return upcoming, nil
}

func (s *AnalyticsService) GetMonthlyRecurringCost(ctx context.Context) (map[models.Currency]float64, error) {
	subscriptions, err := s.subscriptionRepo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	monthlyCosts := make(map[models.Currency]float64)

	for _, sub := range subscriptions {
		// Convert to monthly cost based on period
		monthlyCost := s.calculateMonthlyCost(sub.Cost, sub.PeriodDays)
		monthlyCosts[sub.Currency] += monthlyCost
	}

	return monthlyCosts, nil
}

func (s *AnalyticsService) calculateMonthlyCost(cost float64, periodDays int) float64 {
	// Average days in a month
	avgDaysInMonth := 30.44
	return cost * (avgDaysInMonth / float64(periodDays))
}