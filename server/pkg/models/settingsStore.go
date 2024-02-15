package models

import (
	"database/sql"
	"fmt"
)

type SettingsStore interface {
	Set(settings Settinsg) error
	Get() (Settinsg, error)
}

type SettingsMySqlStore struct {
	DB *sql.DB
}

func NewSettingsStore(db *sql.DB) *SettingsMySqlStore {
	return &SettingsMySqlStore{
		DB: db,
	}
}

func (s SettingsMySqlStore) setSetting(value int, opType int) error {
	sql := "UPDATE settings SET value=? WHERE type=?"
	update, err := s.DB.Query(sql, value, opType)
	if err != nil {
		return err
	}
	defer update.Close()
	return nil
}

func (s SettingsMySqlStore) Set(settings Settinsg) error {
	err := s.setSetting(settings.OpPlusTime, TypeSettingPlus)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpMinusTime, TypeSettingMinus)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpMultTime, TypeSettingMult)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpDivTime, TypeSettingDiv)
	if err != nil {
		return err
	}
	err = s.setSetting(settings.OpAgentTimeOut, TypeSettingAgentTimeout)
	if err != nil {
		return err
	}

	err = s.setSetting(settings.OpAgentDeleteTimeOut, TypeSettingAgentDeleteTimeout)
	if err != nil {
		return err
	}
	return nil
}

func (s SettingsMySqlStore) Get() (Settinsg, error) {
	settings := Settinsg{
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
		case TypeSettingPlus:
			settings.OpPlusTime = value
		case TypeSettingMinus:
			settings.OpMinusTime = value
		case TypeSettingMult:
			settings.OpMultTime = value
		case TypeSettingDiv:
			settings.OpDivTime = value
		case TypeSettingAgentTimeout:
			settings.OpAgentTimeOut = value
		case TypeSettingAgentDeleteTimeout:
			settings.OpAgentDeleteTimeOut = value
		default:
			return settings, fmt.Errorf("not allowed settings type")
		}
	}
	return settings, nil
}
