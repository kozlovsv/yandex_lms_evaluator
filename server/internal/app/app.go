package app

import (
	"log/slog"

	"github.com/kozlovsv/evaluator/server/config"
	grpcapp "github.com/kozlovsv/evaluator/server/internal/app/grpc"
	"github.com/kozlovsv/evaluator/server/internal/checker"
	grpcserv "github.com/kozlovsv/evaluator/server/internal/grpc"
	"github.com/kozlovsv/evaluator/server/internal/storage"
)

type App struct {
	GrpcApp *grpcapp.App
	log     *slog.Logger
	Storage *storage.Storage
	checker *checker.Checker
}

func New(
	log *slog.Logger,
	dbConf config.DBConfig,
	GRPCPort string,
) *App {
	storage, err := storage.New(dbConf, log)

	if err != nil {
		panic(err)
	}

	grpcService := grpcserv.New(log, storage.AgentStore, storage.ExpressionStore, storage.SettingsStore)

	grpcApp := grpcapp.New(log, GRPCPort, grpcService)

	checker := checker.NewChecker(storage.ExpressionStore, storage.AgentStore, storage.SettingsStore, log)

	return &App{
		GrpcApp: grpcApp,
		Storage: storage,
		log:     log,
		checker: checker,
	}
}

// Run runs gRPC server.
func (a *App) Run() {
	a.GrpcApp.MustRun()
	a.checker.Run()
}

// Stop stops gRPC server.
func (a *App) Stop() {
	a.checker.Stop()
	a.GrpcApp.Stop()
	a.Storage.Close()
}
