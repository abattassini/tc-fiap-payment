package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abattassini/tc-fiap-payment/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	app := app.InitializeApp()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("Error while starting app: %v", err)
	}

	<-ctx.Done()

	if err := app.Stop(ctx); err != nil {
		log.Fatalf("Error while stopping app: %v", err)
	}
}
