package storage

import (
	"database/sql"
	"fmt"

	"github.com/kozlovsv/evaluator/api/internal/models"
)

var agentFieldsSet = "`name`, `status`, `last_ping`, `current_op`"

type AgentStore struct {
	DB *sql.DB
}

func NewAgentStore(db *sql.DB) *AgentStore {
	return &AgentStore{
		DB: db,
	}
}

func (s AgentStore) getAgent(row RowScanner) (models.Agent, error) {
	agent := models.Agent{}
	err := row.Scan(&agent.Name, &agent.Status, &agent.LastPing, &agent.CurrentOp)
	return agent, err
}

func (s AgentStore) List() ([]models.Agent, error) {
	var res []models.Agent
	query := fmt.Sprintf("SELECT %s FROM agent ORDER BY name", agentFieldsSet)
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item, err := s.getAgent(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}
