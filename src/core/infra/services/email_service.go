package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type InMemoryEmailService struct{}

func NewInMemoryEmailService() *InMemoryEmailService {
	return &InMemoryEmailService{}
}

func (s *InMemoryEmailService) SendConfirmationEmail(ctx context.Context, userID, eventID string) error {
	if sagaID, ok := ctx.Value("sagaID").(string); ok {
		log.Printf("sPid=%s | email service  evento : %s", sagaID, eventID)
	}
	if rand.Intn(10) == 0 {
		return errors.New("falha ao enviar e-mail: serviço temporariamente indisponível")
	}

	time.Sleep(1 * time.Second) // pequeno delay simulado
	log.Printf("📧 E-mail enviado para o usuário %s sobre o evento %s", userID, eventID)
	fmt.Printf("✉️  Confirmação de compra enviada para o usuário: %s\n", userID)
	return nil
}
