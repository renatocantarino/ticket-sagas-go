package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/renatocantarino/sagas/src/core/domain"
	"github.com/renatocantarino/sagas/src/core/infra/repository"
)

var (
	ErrInvalidQuantity     = errors.New("quantidade inválida")
	ErrInsufficientTickets = errors.New("ingressos insuficientes disponíveis")
	ErrCannotConfirm       = errors.New("ingresso não pode ser confirmado")
	ErrCannotCancel        = errors.New("ingresso não pode ser cancelado")
	ErrAlreadyCancelled    = errors.New("ingresso já está cancelado")
	ErrNotFound            = errors.New("ingresso não encontrado")
)

type InMemoryTicketService struct {
	tickets      map[string]*domain.Ticket
	eventRepo    domain.EventRepository
	domainEvents *repository.InMemoryDomainEventDB

	mu sync.RWMutex
}

func NewInMemoryTicketService(eventRepo domain.EventRepository, domainEvents *repository.InMemoryDomainEventDB) *InMemoryTicketService {
	return &InMemoryTicketService{
		tickets:      make(map[string]*domain.Ticket),
		eventRepo:    eventRepo,
		mu:           sync.RWMutex{},
		domainEvents: domainEvents,
	}
}

func NewTicket(event *domain.Event, userID string, quantity int) (*domain.Ticket, error) {

	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	if event.SoldTickets+quantity > event.MaxTickets {
		return nil, ErrInsufficientTickets
	}

	ticket := &domain.Ticket{
		ID:         fmt.Sprintf("tkt-%d", time.Now().UnixNano()),
		UserID:     userID,
		EventID:    event.ID,
		Status:     domain.TicketReserved,
		ReservedAt: time.Now(),
		Quantity:   quantity,
	}

	ticket.TotalPrice = event.UnitPrice * float64(quantity)
	event.SoldTickets += quantity

	return ticket, nil
}

func (s *InMemoryTicketService) Reserve(ctx context.Context, userID, eventID string, quantity int) (*domain.Ticket, error) {
	sagaID, ok := ctx.Value("saga_id").(string)
	if ok {
		log.Printf("sPid=%s | reserva de ingresso: %s", sagaID, eventID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("evento não encontrado: %w", err)
	}

	if rand.Intn(10) < 3 { // 30% de chance de falha
		return nil, errors.New("erro aleatório na reserva")
	}

	// 3. Validar e criar ticket
	ticket, err := NewTicket(event, userID, quantity)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar ticket: %w", err)
	}

	err = s.eventRepo.Update(event)
	if err != nil {
		return nil, fmt.Errorf("falha ao atualizar evento: %w", err)
	}

	s.tickets[ticket.ID] = ticket

	log.Printf("Reserva realizada: %d ingressos para evento %s. Total: R$ %.2f",
		quantity, eventID, ticket.TotalPrice)

	payload := map[string]interface{}{
		"ticket_id":   ticket.ID,
		"user_id":     userID,
		"event_id":    event.ID,
		"quantity":    quantity,
		"total_price": ticket.TotalPrice,
	}

	err = s.domainEvents.AppendEvent("TicketConfirmed", sagaID, payload)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar event: %w", err)
	}

	return ticket, nil
}

func (s *InMemoryTicketService) Cancel(ctx context.Context, ticketID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sagaID, ok := ctx.Value("sagaID").(string); ok {
		log.Printf("sPid=%s | Cancelando reserva do ingresso: %s", sagaID, ticketID)
	}

	ticket, exists := s.tickets[ticketID]
	if !exists {
		return ErrNotFound
	}

	if ticket.Status == domain.TicketCancelled {
		return ErrAlreadyCancelled
	}

	if ticket.Status != domain.TicketReserved {
		return fmt.Errorf("ingresso não pode ser cancelado no status: %s", ticket.Status)
	}

	event, err := s.eventRepo.FindByID(ticket.EventID)
	if err != nil {
		return fmt.Errorf("erro ao buscar evento: %w", err)
	}

	event.SoldTickets -= ticket.Quantity
	if event.SoldTickets < 0 {
		event.SoldTickets = 0
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
		return ErrNotFound
	}
	ticket.Status = domain.TicketConfirmed
	return nil
}
