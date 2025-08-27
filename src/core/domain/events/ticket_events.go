package domain

import "time"

type TicketEvent struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Occurred  time.Time              `json:"occurredAt"`
	SagaID    string                 `json:"sagaId"`
	Aggregate int                    `json:"aggregate,omitempty"`
}

func (e TicketEvent) EventType() string     { return e.Type }
func (e TicketEvent) OccurredAt() time.Time { return e.Occurred }
