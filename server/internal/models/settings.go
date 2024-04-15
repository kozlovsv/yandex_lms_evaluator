package models

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
