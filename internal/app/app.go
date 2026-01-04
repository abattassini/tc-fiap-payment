package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	paymentController "github.com/abattassini/tc-fiap-payment/internal/payment/controller"
	paymentRepositories "github.com/abattassini/tc-fiap-payment/internal/payment/domain/repositories"
	paymentGateways "github.com/abattassini/tc-fiap-payment/internal/payment/gateways"
	paymentApiController "github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/controller"
	paymentClients "github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/clients"
	paymentGatewaysImpl "github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/gateways"
	paymentPersistence "github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/persistence"
	paymentPresenter "github.com/abattassini/tc-fiap-payment/internal/payment/presenter"
	paymentUseCasesAdd "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/addPayment"
	paymentUseCasesGet "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/getPayment"
	paymentUseCasesGetStatus "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/getPaymentStatus"
	paymentUseCasesHandleWebhook "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/handleWebhook"
	paymentUseCasesUpdate "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/updatePayment"

	"github.com/abattassini/tc-fiap-payment/pkg/rest"
	"github.com/abattassini/tc-fiap-payment/pkg/storage/postgres"
)

func InitializeApp() *fx.App {
	return fx.New(
		fx.Provide(
			postgres.NewPostgresDB,
			fx.Annotate(paymentPersistence.NewPaymentRepositoryImpl, fx.As(new(paymentRepositories.PaymentRepository))),
			fx.Annotate(paymentPresenter.NewPaymentPresenterImpl, fx.As(new(paymentPresenter.PaymentPresenter))),
			fx.Annotate(paymentController.NewPaymentControllerImpl, fx.As(new(paymentController.PaymentController))),
			fx.Annotate(paymentController.NewPaymentWebhookControllerImpl, fx.As(new(paymentController.PaymentWebhookController))),
			fx.Annotate(paymentUseCasesAdd.NewAddPaymentUseCaseImpl, fx.As(new(paymentUseCasesAdd.AddPaymentUseCase))),
			fx.Annotate(paymentUseCasesGet.NewGetPaymentUseCaseImpl, fx.As(new(paymentUseCasesGet.GetPaymentUseCase))),
			fx.Annotate(paymentUseCasesGetStatus.NewGetPaymentStatusUseCaseImpl, fx.As(new(paymentUseCasesGetStatus.GetPaymentStatusUseCase))),
			fx.Annotate(paymentUseCasesUpdate.NewUpdatePaymentUseCaseImpl, fx.As(new(paymentUseCasesUpdate.UpdatePaymentUseCase))),
			fx.Annotate(paymentUseCasesHandleWebhook.NewHandleWebhookUseCaseImpl, fx.As(new(paymentUseCasesHandleWebhook.HandleWebhookUseCase))),
			func() (paymentGateways.MercadoPagoGateway, error) {
				return paymentGatewaysImpl.NewMercadoPagoGatewayImpl()
			},
			func() rest.HTTPClient {
				return &http.Client{}
			},
			func(httpClient rest.HTTPClient) paymentClients.OrderClient {
				return paymentClients.NewOrderClient(httpClient)
			},
			chi.NewRouter,
			func(
				paymentController paymentController.PaymentController,
				paymentWebhookController paymentController.PaymentWebhookController) []rest.Controller {
				return []rest.Controller{
					paymentApiController.NewPaymentApiController(paymentController),
					paymentApiController.NewPaymentWebhookApiController(paymentWebhookController),
				}
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
