package domain

import (
	"context"
	"time"
)

type Payment struct {
	ID     string
	UserID string
	Amount float64
	Status PaymentStatus
	PaidAt time.Time
}

type PaymentStatus string

const (
	PaymentProcessed PaymentStatus = "processed"
	PaymentRefunded  PaymentStatus = "refunded"
	PaymentFailed    PaymentStatus = "failed"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, userID string, amount float64) (*Payment, error)
	RefundPayment(ctx context.Context, paymentID string) error
}
