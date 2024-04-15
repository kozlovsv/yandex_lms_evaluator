package storage

import (
	"database/sql"
	"log"
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

func (s AgentStore) Add(agentName string) error {
	//Проверяем может уже есть агент
	name := ""
	err := s.DB.QueryRow("SELECT name FROM agent WHERE name = ? LIMIT 1", agentName).Scan(&name)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if name != "" {
		return nil
	}

	insert, err := s.DB.Query("INSERT INTO agent (name, status, last_ping) VALUES (?, 0, CURRENT_TIMESTAMP())", agentName)

	if err != nil {
		return err
	}
	defer insert.Close()
	return nil
}

func (s AgentStore) SetAvailable(name string) error {
	_, err := s.DB.Exec("UPDATE agent SET status = 0, last_ping = CURRENT_TIMESTAMP() WHERE name = ?", name)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s AgentStore) SetCurrentOp(name string, currentOp string) error {
	_, err := s.DB.Exec("UPDATE agent SET status = 0, last_ping = CURRENT_TIMESTAMP(), current_op = ? WHERE name = ?", currentOp, name)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s AgentStore) SetNotAvailable(timeout int) error {
	stm, err := s.DB.Exec("UPDATE agent SET status = 1 WHERE status = 0 AND last_ping <= DATE_SUB(NOW(), INTERVAL ? * 1000 MICROSECOND)", timeout)
	if err == nil {
		if cnt, _ := stm.RowsAffected(); cnt > 0 {
			log.Println("[INFO] Деактивировано агентов:", cnt, timeout)
		}
	}
	return err
}

func (s AgentStore) DeleteNotAvailable(timeout int) error {
	stm, err := s.DB.Exec("DELETE FROM agent WHERE status = 1 AND last_ping <= DATE_SUB(NOW(), INTERVAL ? * 1000 MICROSECOND)", timeout)
	if err == nil {
		if cnt, _ := stm.RowsAffected(); cnt > 0 {
			log.Println("[INFO] Удалено агентов:", cnt, timeout)
		}
	}
	return err
}
