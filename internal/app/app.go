package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	paymentRepositories "github.com/abattassini/tc-fiap-payment/internal/payment/domain/repositories"
	paymentPersistence "github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/persistence"

	"github.com/abattassini/tc-fiap-payment/pkg/rest"
	"github.com/abattassini/tc-fiap-payment/pkg/storage/postgres"
)

func InitializeApp() *fx.App {
	return fx.New(
		fx.Provide(
			postgres.NewPostgresDB,
			fx.Annotate(paymentPersistence.NewPaymentRepositoryImpl, fx.As(new(paymentRepositories.PaymentRepository))),
			chi.NewRouter,
			func() []rest.Controller {
				return []rest.Controller{}
			},
		),
		fx.Invoke(registerRoutes),
		fx.Invoke(startHTTPServer),
	)
}

func registerRoutes(r *chi.Mux, controllers []rest.Controller) {
	r.Use(middleware.Logger)

	for _, controller := range controllers {
		controller.RegisterRoutes(r)
	}
}

func startHTTPServer(lc fx.Lifecycle, r *chi.Mux) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Starting HTTP server on :%s", port)
				if err := http.ListenAndServe(":"+port, r); err != nil {
					log.Fatalf("Failed to start HTTP server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server gracefully")
			return nil
		},
	})
}
