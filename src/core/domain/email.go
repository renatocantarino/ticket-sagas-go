package domain

import "context"

type EmailService interface {
	SendConfirmationEmail(ctx context.Context, userID, eventID string) error
}
