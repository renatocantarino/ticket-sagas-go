package application

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/renatocantarino/sagas/src/core/domain"
)

type SagaStep func() error
type Compensation func()

const (
	sagaIDKey string = "saga_id"
)

type SagaStepWithCompensation struct {
	Step         SagaStep
	Compensation Compensation
}

type TicketPurchaseSaga struct {
	ticketSvc  domain.TicketService
	paymentSvc domain.PaymentService
	emailSvc   domain.EmailService
}

func NewTicketPurchaseSaga(
	ticketSvc domain.TicketService,
	paymentSvc domain.PaymentService,
	emailSvc domain.EmailService,
) *TicketPurchaseSaga {
	return &TicketPurchaseSaga{
		ticketSvc:  ticketSvc,
		paymentSvc: paymentSvc,
		emailSvc:   emailSvc,
	}
}

func (s *TicketPurchaseSaga) Handler(userID, eventID string, quantity int) error {

	sagaID := uuid.New().String()
	ctx := context.WithValue(context.Background(), sagaIDKey, sagaID)

	log.Printf("Saga iniciada [saga_id=%s] | Usuário: %s, Evento: %s", sagaID, userID, eventID)

	var ticket *domain.Ticket
	var payment *domain.Payment

	steps := []SagaStepWithCompensation{
		{
			Step: func() error {
				var err error
				ticket, err = s.ticketSvc.Reserve(ctx, userID, eventID, quantity)
				if err != nil {
					log.Printf("Falha em Step Reserve [saga_id=%s]: %v", sagaID, err)
					return err
				}
				log.Printf("Ingresso reservado [saga_id=%s] | TicketID: %s", sagaID, ticket.ID)
				return nil
			},
			Compensation: func() {
				if ticket != nil {
					_ = s.ticketSvc.Cancel(ctx, ticket.ID)
					log.Printf("Reserva cancelada [saga_id=%s] | TicketID: %s", sagaID, ticket.ID)
				}
			},
		},
		{
			Step: func() error {
				var err error
				payment, err = s.paymentSvc.ProcessPayment(ctx, userID, ticket.TotalPrice)
				if err != nil {
					return err
				}
				log.Printf("Pagamento processado: %s", payment.ID)
				return nil
			},
			Compensation: func() {
				if payment != nil {
					_ = s.paymentSvc.RefundPayment(ctx, payment.ID)
					log.Printf("Pagamento reembolsado: %s", payment.ID)
				}
			},
		},
		{
			Step: func() error {
				return s.emailSvc.SendConfirmationEmail(ctx, userID, eventID)
			},
			Compensation: func() {
				log.Printf("⚠️  Não é possível desfazer envio de e-mail para %s", userID)
			},
		},
	}

	var compensations []Compensation

	for i, step := range steps {
		if err := step.Step(); err != nil {
			log.Printf("❌ Falha no passo %d: %v", i+1, err)
			for j := len(compensations) - 1; j >= 0; j-- {
				compensations[j]()
			}
			return fmt.Errorf("saga falhou no passo %d: %w", i+1, err)
		}
		compensations = append(compensations, step.Compensation)
	}

	log.Println("✅ Saga concluída com sucesso!")
	return nil
}
