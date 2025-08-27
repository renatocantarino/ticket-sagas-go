package infra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/renatocantarino/sagas/src/core/domain"
)

type InMemoryPaymentService struct {
	payments map[string]*domain.Payment
}

func NewInMemoryPaymentService() *InMemoryPaymentService {
	return &InMemoryPaymentService{
		payments: make(map[string]*domain.Payment),
	}
}

func (s *InMemoryPaymentService) ProcessPayment(ctx context.Context, userID string, amount float64) (*domain.Payment, error) {

	if amount <= 0 {
		return nil, errors.New("valor do pagamento deve ser maior que zero")
	}

	// Simula falha aleatória (20% de chance)
	if rand.Intn(10) < 2 {
		return nil, errors.New("pagamento falhou: erro no gateway")
	}

	// Simula recusa por limite
	if amount > 1000 {
		return nil, errors.New("pagamento recusado: limite de valor excedido")
	}

	payment := &domain.Payment{
		ID:     fmt.Sprintf("pay-%d", time.Now().UnixNano()),
		UserID: userID,
		Amount: amount,
		Status: domain.PaymentProcessed,
		PaidAt: time.Now(),
	}

	if sagaID, ok := ctx.Value("sagaID").(string); ok {
		log.Printf("sPid=%s | pagando ingresso: %s", sagaID, payment.ID)
	}

	s.payments[payment.ID] = payment

	return payment, nil
}

func (s *InMemoryPaymentService) RefundPayment(ctx context.Context, paymentID string) error {

	if sagaID, ok := ctx.Value("sagaID").(string); ok {
		log.Printf("sPid=%s | RefundPayment  ingresso: %s", sagaID, paymentID)
	}

	payment, exists := s.payments[paymentID]
	if !exists {
		return errors.New("pagamento não encontrado")
	}

	if payment.Status == domain.PaymentRefunded {
		return nil
	}

	payment.Status = domain.PaymentRefunded
	return nil
}
