package repository

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type DomainEvents struct {
	SagaID    string    `json:"saga_id"`
	Type      string    `json:"type"`
	Payload   []byte    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	Version   int       `json:"version"`
}

type InMemoryDomainEventDB struct {
	mu    sync.RWMutex
	store map[string][]DomainEvents // sagaID -> []Event
}

func NewInMemoryEventDB() *InMemoryDomainEventDB {
	return &InMemoryDomainEventDB{
		store: make(map[string][]DomainEvents),
	}
}

func BuildDomainEvent(eventType, sagaID string, payload interface{}) (DomainEvents, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return DomainEvents{}, err
	}

	return DomainEvents{
		SagaID:    sagaID,
		Type:      eventType,
		Payload:   payloadBytes,
		Timestamp: time.Now().UTC(),
		Version:   1,
	}, nil
}

func (db *InMemoryDomainEventDB) AppendEvent(eventType, SagaId string, payload any) error {

	log.Printf("sPid=%s | criando eventStore do tipo:%s", SagaId, eventType)

	db.mu.Lock()
	defer db.mu.Unlock()

	dm, err := BuildDomainEvent(eventType, SagaId, payload)
	if err != nil {

		return err
	}

	events := db.store[dm.SagaID]
	dm.Version = len(events) + 1
	db.store[dm.SagaID] = append(events, dm)

	return nil
}

func (db *InMemoryDomainEventDB) GetEvents(sagaID string) []DomainEvents {
	db.mu.RLock()
	defer db.mu.RUnlock()

	events, ok := db.store[sagaID]
	if !ok {
		return []DomainEvents{}
	}

	copied := make([]DomainEvents, len(events))
	copy(copied, events)
	return copied
}

func (db *InMemoryDomainEventDB) HasSaga(sagaID string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, exists := db.store[sagaID]
	return exists
}

func (db *InMemoryDomainEventDB) DeleteSaga(sagaID string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.store, sagaID)
}
