package main

import (
	"log"

	application "github.com/renatocantarino/sagas/src/core/app"
	"github.com/renatocantarino/sagas/src/core/infra"
)

func main() {
	ticketSvc := infra.NewInMemoryEventService()
	paymentSvc := infra.NewInMemoryPaymentService()
	emailSvc := infra.NewInMemoryEmailService()
	eventRepo := infra.NewInMemoryEventRepository()

	//saga creation
	saga := application.NewTicketPurchaseSaga(ticketSvc, paymentSvc, emailSvc, eventRepo)

	userID := "user-789"
	eventID := "event-001"
	quantity := 2

	log.Printf("Iniciando compra de %d ingresso(s) para %s no evento %s para o usuario (%s)",
		quantity,
		userID,
		eventID,
	)
	err := saga.Handler(userID, eventID, quantity)
	if err != nil {
		log.Fatalf("❌ Compra falhou: %v", err)
	}

	log.Println("🎉 Compra concluída com sucesso!")

}
