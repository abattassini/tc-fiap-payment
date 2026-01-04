package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	RegisterRoutes(r chi.Router)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
