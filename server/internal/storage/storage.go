package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kozlovsv/evaluator/server/config"
)

type Storage struct {
	db              *sql.DB
	log             *slog.Logger
	AgentStore      *AgentStore
	ExpressionStore *ExpressionStore
	SettingsStore   *SettingsStore
}

func New(dbConf config.DBConfig, log *slog.Logger) (*Storage, error) {
	db, err := sql.Open("mysql", dbConf.GetDSNString())

	if err != nil {
		return nil, fmt.Errorf("%s: %w", "storage.New", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "storage.New", err)
	}
	log.Info("Storage SERVER MYSQL Startes")

	return &Storage{
		db:              db,
		log:             log,
		ExpressionStore: NewExpressionStore(db),
		AgentStore:      NewAgentStore(db),
		SettingsStore:   NewSettingsStore(db),
	}, nil
}

func (s *Storage) Close() error {
	s.log.Info("Storage SERVER MYSQL Stoped")
	return s.db.Close()
}
