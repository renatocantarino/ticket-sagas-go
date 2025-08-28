package domain

type Event struct {
	ID          string
	Name        string
	Date        string
	UnitPrice   float64
	MaxTickets  int
	SoldTickets int
}

type EventRepository interface {
	FindByID(eventID string) (*Event, error)
	Update(event *Event) error
}
