package infra

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	domain "github.com/renatocantarino/sagas/src/core/domain/events"
)

type EventStoreDb struct {
	db *sql.DB
}

func NewSQLEventStore(dbPath string) (*EventStoreDb, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	store := &EventStoreDb{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *EventStoreDb) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS domain_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_type TEXT NOT NULL,
		payload JSON NOT NULL,
		occurred_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_event_type ON domain_events(event_type);
	CREATE INDEX IF NOT EXISTS idx_occurred_at ON domain_events(occurred_at);
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *EventStoreDb) Save(event domain.TicketEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO domain_events (event_type, payload, occurred_at) VALUES (?, ?, ?)",
		event.EventType(),
		string(payload),
		event.OccurredAt().Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		log.Printf("❌ Falha ao salvar evento %s: %v", event.EventType(), err)
	} else {
		log.Printf("✅ Evento salvo: %s", event.EventType())
	}

	return err
}

func (s *EventStoreDb) GetAllEventsByTicketId(ticketId int) ([]domain.TicketEvent, error) {

	query := `
		SELECT event_type, payload, occurred_at, saga_id 
		FROM domain_events 
		WHERE Aggregate = ?
		ORDER BY occurred_at
	`

	rows, err := s.db.Query(query, ticketId)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar banco: %w", err)
	}
	defer rows.Close()

	var events []domain.TicketEvent
	for rows.Next() {
		var eventType, payload, occurredAt, sagaID string
		if err := rows.Scan(&eventType, &payload, &occurredAt, &sagaID); err != nil {
			return nil, fmt.Errorf("falha ao ler linha: %w", err)
		}

		// Converter timestamp para time.Time
		occurred, err := time.Parse("2006-01-02 15:04:05", occurredAt)
		if err != nil {
			return nil, fmt.Errorf("timestamp inválido: %s, erro: %w", occurredAt, err)
		}
		var payloadData map[string]interface{}
		if err := json.Unmarshal([]byte(payload), &payloadData); err != nil {
			return nil, fmt.Errorf("falha ao desserializar payload do evento %s: %w", eventType, err)
		}

		event := domain.TicketEvent{
			Type:      eventType,
			Payload:   payloadData,
			Occurred:  occurred,
			SagaID:    sagaID,
			Aggregate: ticketId,
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração dos resultados: %w", err)
	}

	log.Printf("✅ Carregados %d eventos do banco.", len(events))
	return events, nil
}
