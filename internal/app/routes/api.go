package routes

import (
	"net/http"

	"github.com/J0es1ick/shortli/internal/app/handlers/urlHandlers"
	"github.com/J0es1ick/shortli/internal/config"
	"github.com/J0es1ick/shortli/internal/repository"
)

func SetupRoutes(cfg *config.Config, urlRepository *repository.UrlRepository) http.Handler {
	mux := http.NewServeMux()

	urlHandler := urlHandlers.NewHandler(cfg, urlRepository)
    
    mux.HandleFunc("GET /", urlHandler.Home)
    mux.HandleFunc("POST /api/shorten", urlHandler.Shorten)
    mux.HandleFunc("GET /api/stats/{shortCode}", urlHandler.Stats)
    mux.HandleFunc("GET /{shortCode}", urlHandler.Redirect)
    mux.HandleFunc("DELETE /urls/{shortCode}", urlHandler.Delete)

	return mux
}
