package domain

import (
	"context"
	"time"
)

type TicketStatus string

const (
	TicketReserved  TicketStatus = "reserved"
	TicketCancelled TicketStatus = "cancelled"
	TicketConfirmed TicketStatus = "confirmed"
)

type Ticket struct {
	ID         string
	EventID    string
	UserID     string
	Status     TicketStatus
	ReservedAt time.Time
	Quantity   int
	TotalPrice float64
}

type TicketService interface {
	Reserve(ctx context.Context, userId, eventId string, qtd int) (*Ticket, error)
	Cancel(ctx context.Context, tickeId string) error
	Confirm(ctx context.Context, tickeId string) error
}
