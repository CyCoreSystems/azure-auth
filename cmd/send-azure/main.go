package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/CyCoreSystems/azure-auth/pkg/config"
	"github.com/CyCoreSystems/azure-auth/pkg/smtp"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	if err := smtp.Send(ctx, cfg, "smtp.office365.com:587", os.Args[1:], os.Stdin); err != nil {
		log.Fatalln("failed to send message:", err)
	}
}
