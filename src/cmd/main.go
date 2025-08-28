package main

import (
	"log"

	application "github.com/renatocantarino/sagas/src/core/app"
	"github.com/renatocantarino/sagas/src/core/infra"
)

func main() {

	eventRepo := infra.NewInMemoryEventRepository()

	ticketSvc := infra.NewInMemoryTicketService(eventRepo)
	paymentSvc := infra.NewInMemoryPaymentService()
	emailSvc := infra.NewInMemoryEmailService()

	//saga creation
	saga := application.NewTicketPurchaseSaga(ticketSvc, paymentSvc, emailSvc)

	userID := "user-789"
	eventID := "event-001"
	quantity := 2

	log.Printf("Iniciando compra de %d ingresso(s) para o evento %s (usuário: %s)",
		quantity,
		eventID,
		userID,
	)

	err := saga.Handler(userID, eventID, quantity)
	if err != nil {
		log.Fatalf("❌ Compra falhou: %v", err)
	}

	log.Println("🎉 Compra concluída com sucesso!")

}
