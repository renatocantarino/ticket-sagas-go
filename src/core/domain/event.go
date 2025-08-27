package domain

type Event struct {
	ID          string
	Name        string
	Date        string
	Price       float64
	MaxTickets  int
	SoldTickets int
}

type EventRepository interface {
	FindByID(eventID string) (*Event, error)
}
