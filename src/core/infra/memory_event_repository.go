package infra

import (
	"errors"
	"sync"

	"github.com/renatocantarino/sagas/src/core/domain"
)

type InMemoryEventRepository struct {
	mu     sync.RWMutex
	events map[string]*domain.Event
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	repo := &InMemoryEventRepository{
		events: make(map[string]*domain.Event),
	}

	// Adiciona eventos de exemplo
	repo.events["event-001"] = &domain.Event{
		ID:          "event-001",
		Name:        "Concerto de Rock",
		Date:        "2025-06-15T20:00:00Z",
		Price:       150.0,
		MaxTickets:  100,
		SoldTickets: 0,
	}

	repo.events["event-002"] = &domain.Event{
		ID:          "event-002",
		Name:        "Workshop de Go",
		Date:        "2025-05-20T09:00:00Z",
		Price:       80.0,
		MaxTickets:  50,
		SoldTickets: 5,
	}

	return repo
}

func (r *InMemoryEventRepository) FindByID(eventID string) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[eventID]
	if !exists {
		return nil, errors.New("evento não encontrado")
	}

	// Retornar cópia para evitar mutação acidental
	cloned := *event
	return &cloned, nil
}

// Optional: método para adicionar eventos (útil em testes)
func (r *InMemoryEventRepository) AddEvent(event *domain.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[event.ID] = event
}
