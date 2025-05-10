package routes

import (
	"net/http"

	"github.com/J0es1ick/shortli/internal/app/handlers/urlHandlers"
	"github.com/J0es1ick/shortli/internal/config"
	"github.com/J0es1ick/shortli/internal/repository"
)

func SetupRoutes(cfg *config.Config, urlRepository *repository.UrlRepository) http.Handler {
	mux := http.NewServeMux()
	handler := urlHandlers.NewHandler(cfg, urlRepository)

	mux.HandleFunc("GET /", handler.Home)
	mux.HandleFunc("POST /api/shorten", handler.Shorten)
	mux.HandleFunc("GET /api/stats/{shortCode}", handler.Stats)
	mux.HandleFunc("GET /{shortCode}", handler.Redirect)

	return mux
}
