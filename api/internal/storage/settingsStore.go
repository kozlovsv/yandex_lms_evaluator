package storage

import (
	"database/sql"
	"fmt"

	"github.com/kozlovsv/evaluator/api/internal/models"
)

type SettingsStore struct {
	DB *sql.DB
}

func NewSettingsStore(db *sql.DB) *SettingsStore {
	return &SettingsStore{
		DB: db,
	}
}

func (s SettingsStore) setSetting(value int, opType int) error {
	sql := "UPDATE settings SET value=? WHERE type=?"
	update, err := s.DB.Query(sql, value, opType)
	if err != nil {
		return err
	}
	defer update.Close()
	return nil
}

func (s SettingsStore) Set(settings models.Settinsg) error {
	err := s.setSetting(settings.OpPlusTime, models.TypeSettingPlus)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpMinusTime, models.TypeSettingMinus)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpMultTime, models.TypeSettingMult)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpDivTime, models.TypeSettingDiv)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpAgentTimeOut, models.TypeSettingAgentTimeout)
	if err != nil {
		return err
	}

	err = s.setSetting(settings.OpAgentDeleteTimeOut, models.TypeSettingAgentDeleteTimeout)
	if err != nil {
		return err
	}
	return nil
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
