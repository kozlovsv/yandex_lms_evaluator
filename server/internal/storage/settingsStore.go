package storage

import (
	"database/sql"
	"fmt"

	"github.com/kozlovsv/evaluator/server/internal/models"
)

type SettingsStore struct {
	DB *sql.DB
}

func NewSettingsStore(db *sql.DB) *SettingsStore {
	return &SettingsStore{
		DB: db,
	}
}

func (s SettingsStore) Get() (models.Settinsg, error) {
	settings := models.Settinsg{
		OpPlusTime:           300,
		OpMinusTime:          300,
		OpMultTime:           300,
		OpDivTime:            300,
		OpAgentTimeOut:       300,
		OpAgentDeleteTimeOut: 300,
	}
	rows, err := s.DB.Query("SELECT type + 0 as id, value FROM settings")
	if err != nil {
		return settings, err
	}
	defer rows.Close()
	var (
		id    int
		value int
	)
	for rows.Next() {
		err := rows.Scan(&id, &value)
		if err != nil {
			return settings, err
		}
		switch id {
		case models.TypeSettingPlus:
			settings.OpPlusTime = value
		case models.TypeSettingMinus:
			settings.OpMinusTime = value
		case models.TypeSettingMult:
			settings.OpMultTime = value
		case models.TypeSettingDiv:
			settings.OpDivTime = value
		case models.TypeSettingAgentTimeout:
			settings.OpAgentTimeOut = value
		case models.TypeSettingAgentDeleteTimeout:
			settings.OpAgentDeleteTimeOut = value
		default:
			return settings, fmt.Errorf("not allowed settings type")
		}
	}
	return settings, nil
}
