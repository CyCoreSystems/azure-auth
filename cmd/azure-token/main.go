package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/CyCoreSystems/azure-auth/pkg/config"
	"github.com/CyCoreSystems/azure-auth/pkg/token"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	tok, err := token.Get(ctx, cfg)
	if err != nil {
		log.Fatalln("failed to get token:", err)
	}

	fmt.Println(tok.AccessToken)
}
