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

type InMemoryTicketService struct {
	tickets map[string]*domain.Ticket
}

func NewInMemoryEventService() *InMemoryTicketService {
	return &InMemoryTicketService{
		tickets: make(map[string]*domain.Ticket),
	}
}

func (s *InMemoryTicketService) Reserve(ctx context.Context, userID, eventID string, qtd int) (*domain.Ticket, error) {
	if sagaID, ok := ctx.Value("sagaID").(string); ok {
		log.Printf("sPid=%s | reserva de ingresso: %s", sagaID, eventID)
	}

	if rand.Intn(10) < 3 { // 30% de chance de falha
		return nil, errors.New("ingressos esgotados")
	}

	ticket := &domain.Ticket{
		ID:         fmt.Sprintf("tkt-%d", time.Now().UnixNano()),
		UserID:     userID,
		EventID:    eventID,
		Status:     domain.TicketReserved,
		ReservedAt: time.Now(),
	}

	s.tickets[ticket.ID] = ticket
	return ticket, nil
}

func (s *InMemoryTicketService) Cancel(ctx context.Context, ticketID string) error {

	if sagaID, ok := ctx.Value("sagaID").(string); ok {
		log.Printf("sPid=%s | Cancelando reserva do ingresso: %s", sagaID, ticketID)
	}

	ticket, exists := s.tickets[ticketID]
	if !exists {
		return errors.New("ingresso não encontrado")
	}
	ticket.Status = domain.TicketCancelled
	return nil
}

func (s *InMemoryTicketService) Confirm(ctx context.Context, ticketID string) error {
	if sagaID, ok := ctx.Value("sagaID").(string); ok {
		log.Printf("sPid=%s | Confirm reserva do ingresso: %s", sagaID, ticketID)
	}

	ticket, exists := s.tickets[ticketID]
	if !exists {
		return errors.New("ingresso não encontrado")
	}
	ticket.Status = domain.TicketConfirmed
	return nil
}
