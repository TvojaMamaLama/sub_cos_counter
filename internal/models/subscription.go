package models

import (
	"time"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyRUB Currency = "RUB"
)

type Category string

const (
	CategoryEntertainment Category = "entertainment"
	CategoryWork          Category = "work"
	CategoryEducation     Category = "education"
	CategoryHome          Category = "home"
	CategoryOther         Category = "other"
)

type Subscription struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Cost        Money     `json:"cost"`
	Currency    Currency  `json:"currency"`
	PeriodDays  int       `json:"period_days"`
	NextPayment time.Time `json:"next_payment"`
	Category    Category  `json:"category"`
	AutoRenewal bool      `json:"auto_renewal"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	Name        string    `json:"name"`
	Cost        Money     `json:"cost"`
	Currency    Currency  `json:"currency"`
	PeriodDays  int       `json:"period_days"`
	NextPayment time.Time `json:"next_payment"`
	Category    Category  `json:"category"`
	AutoRenewal bool      `json:"auto_renewal"`
}

func (s *Subscription) IsPaymentDue() bool {
	return time.Now().After(s.NextPayment)
}

func (s *Subscription) UpdateNextPayment() {
	s.NextPayment = s.NextPayment.AddDate(0, 0, s.PeriodDays)
	s.UpdatedAt = time.Now()
}
