package httpapp

import (
	"log/slog"
	"net/http"

	"github.com/kozlovsv/evaluator/sso/config"
	"github.com/kozlovsv/evaluator/sso/internal/handlers"
	auth "github.com/kozlovsv/evaluator/sso/internal/services"
)

type App struct {
	log        *slog.Logger
	auth       auth.AuthIntfs
	httpConfig config.HttpConfig
}

func New(log *slog.Logger, auth auth.AuthIntfs, httpConfig config.HttpConfig) *App {
	return &App{
		log:        log,
		httpConfig: httpConfig,
		auth:       auth,
	}
}

// Run runs gRPC server.
func (a *App) Run() {
	const op = "http.app.Run"

	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.Handle("/login", handlers.NewLoginHandler(a.auth, a.log))
	mux.Handle("/register", handlers.NewRegisterHandler(a.auth, a.log))
	a.log.With(slog.String("op", op)).
		Info("Starting SSO HTTP server ON ", slog.String("port", a.httpConfig.Port))
	http.ListenAndServe(a.httpConfig.Host+":"+a.httpConfig.Port, mux)
}
