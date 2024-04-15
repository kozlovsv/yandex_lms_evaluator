package httpapp

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kozlovsv/evaluator/api/config"
	"github.com/kozlovsv/evaluator/api/internal/handlers"
	auth "github.com/kozlovsv/evaluator/api/internal/services"
	"github.com/kozlovsv/evaluator/api/internal/storage"
)

type App struct {
	log        *slog.Logger
	auth       auth.AuthIntfs
	httpConfig config.HttpConfig
	storage    *storage.Storage
}

func New(log *slog.Logger, storage *storage.Storage, auth auth.AuthIntfs, httpConfig config.HttpConfig) *App {
	return &App{
		log:        log,
		httpConfig: httpConfig,
		auth:       auth,
		storage:    storage,
	}
}

// // extractBearerToken extracts auth token from Authorization header.
func ExtractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func JWTMiddleware(auth auth.AuthIntfs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := ExtractBearerToken(r)
		if bearerToken == "" {
			http.Error(w, "Invalid token format", http.StatusBadRequest)
			return
		}
		userId, err := auth.ParseToken(bearerToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.ContextJWTKey, bearerToken)
		ctx = context.WithValue(ctx, handlers.ContextUserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Run runs gRPC server.
func (a *App) Run() {
	const op = "http.app.Run"

	mux := http.NewServeMux()

	expressionHandler := handlers.NewExpressionHandler(a.storage.ExpressionStore)

	// Register the routes and handlers
	mux.HandleFunc("/", handlers.HomeHandler)

	mux.Handle("/expressions", JWTMiddleware(a.auth, expressionHandler))
	mux.Handle("/expressions/", JWTMiddleware(a.auth, expressionHandler))
	mux.Handle("/settings", JWTMiddleware(a.auth, handlers.NewSettingsHandler(a.storage.SettingsStore)))
	mux.Handle("/agents", JWTMiddleware(a.auth, handlers.NewAgentHandler(a.storage.AgentStore)))

	a.log.With(slog.String("op", op)).
		Info("Starting API HTTP server ON ", slog.String("port", a.httpConfig.Port))
	http.ListenAndServe(a.httpConfig.Host+":"+a.httpConfig.Port, mux)
	a.log.With(slog.String("op", op)).
		Info("Auth HTTP Server Stoped")
}
