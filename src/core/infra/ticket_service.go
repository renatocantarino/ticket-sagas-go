package infra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/renatocantarino/sagas/src/core/domain"
)

type InMemoryTicketService struct {
	tickets   map[string]*domain.Ticket
	eventRepo domain.EventRepository
	mu        sync.RWMutex
}

func NewInMemoryTicketService(eventRepo domain.EventRepository) *InMemoryTicketService {
	return &InMemoryTicketService{
		tickets:   make(map[string]*domain.Ticket),
		eventRepo: eventRepo,
		mu:        sync.RWMutex{},
	}
}

func NewTicket(event *domain.Event, userID string, quantity int) (*domain.Ticket, error) {

	if quantity <= 0 {
		return nil, errors.New("quantidade inválida")
	}

	if event.SoldTickets+quantity > event.MaxTickets {
		return nil, errors.New("ingressos insuficientes disponíveis")
	}

	ticket := &domain.Ticket{
		ID:         fmt.Sprintf("tkt-%d", time.Now().UnixNano()),
		UserID:     userID,
		EventID:    event.ID,
		Status:     domain.TicketReserved,
		ReservedAt: time.Now(),
		Quantity:   quantity, // ✅ Não esqueça de setar a quantidade!
	}

	ticket.TotalPrice = event.UnitPrice * float64(quantity)

	return ticket, nil
}

func (s *InMemoryTicketService) Reserve(ctx context.Context, userID, eventID string, quantity int) (*domain.Ticket, error) {
	if sagaID, ok := ctx.Value("sagaID").(string); ok {
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

	event.SoldTickets += quantity
	err = s.eventRepo.Update(event)
	if err != nil {
		return nil, fmt.Errorf("falha ao atualizar evento: %w", err)
	}

	s.tickets[ticket.ID] = ticket

	log.Printf("Reserva realizada: %d ingressos para evento %s. Total: R$ %.2f",
		quantity, eventID, ticket.TotalPrice)

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
		return errors.New("ingresso não encontrado")
	}

	if ticket.Status == domain.TicketCancelled {
		return errors.New("ingresso já está cancelado")
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
		return errors.New("ingresso não encontrado")
	}
	ticket.Status = domain.TicketConfirmed
	return nil
}
