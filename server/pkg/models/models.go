package models

import "time"

type Expression struct {
	Value          string    `json:"value"`
	IdempotencyKey string    `json:"idempotency_key"`
	Id             int       `json:"id"`
	Status         int       `json:"status"`
	Result         string    `json:"result"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatesAt      time.Time `json:"created_at"`
}

const TypeSettingPlus = 1
const TypeSettingMinus = 2
const TypeSettingMult = 3
const TypeSettingDiv = 4
const TypeSettingAgentTimeout = 5
const TypeSettingAgentDeleteTimeout = 6

type Settinsg struct {
	OpPlusTime           int `json:"op_plus"`
	OpMinusTime          int `json:"op_minus"`
	OpMultTime           int `json:"op_mult"`
	OpDivTime            int `json:"op_div"`
	OpAgentTimeOut       int `json:"op_agent_timeout"`
	OpAgentDeleteTimeOut int `json:"op_agent_deletetimeout"`
}

type Task struct {
	Expression Expression
	Settings   Settinsg
}

type Agent struct {
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	LastPing  time.Time `json:"last_ping"`
	CurrentOp string    `json:"current_op"`
}
