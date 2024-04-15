package app

import (
	"log/slog"

	"github.com/kozlovsv/evaluator/api/config"
	httpapp "github.com/kozlovsv/evaluator/api/internal/app/http"
	auth "github.com/kozlovsv/evaluator/api/internal/services"
	"github.com/kozlovsv/evaluator/api/internal/storage"
)

type App struct {
	HttpServer *httpapp.App
	Storage    *storage.Storage
}

func New(
	log *slog.Logger,
	httpConfig config.HttpConfig,
	dbConf config.DBConfig,
	jwtConfig config.JWTConfig,
) *App {
	storage, err := storage.New(dbConf, log)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, jwtConfig)

	httpApp := httpapp.New(log, storage, authService, httpConfig)

	return &App{
		HttpServer: httpApp,
		Storage:    storage,
	}
}

// Run runs gRPC server.
func (a *App) Run() {
	a.HttpServer.Run()
}

// Stop stops gRPC server.
func (a *App) Stop() {
	a.Storage.Close()
}
