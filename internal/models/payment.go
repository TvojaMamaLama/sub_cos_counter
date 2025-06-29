package models

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID             int           `json:"id"`
	SubscriptionID int           `json:"subscription_id"`
	Amount         Money         `json:"amount"`
	Currency       Currency      `json:"currency"`
	PaidAt         time.Time     `json:"paid_at"`
	Status         PaymentStatus `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
}

type CreatePaymentRequest struct {
	SubscriptionID int           `json:"subscription_id"`
	Amount         Money         `json:"amount"`
	Currency       Currency      `json:"currency"`
	Status         PaymentStatus `json:"status"`
}

type PaymentSummary struct {
	Currency    Currency `json:"currency"`
	TotalAmount Money    `json:"total_amount"`
	Count       int      `json:"count"`
}

type MonthlyExpense struct {
	Month    time.Time        `json:"month"`
	Payments []PaymentSummary `json:"payments"`
}
