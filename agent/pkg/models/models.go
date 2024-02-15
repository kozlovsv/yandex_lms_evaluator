package models

type Expression struct {
	Value string `json:"value"`
	Id    int    `json:"id"`
}

type Settinsg struct {
	OpPlusTime     int `json:"op_plus"`
	OpMinusTime    int `json:"op_minus"`
	OpMultTime     int `json:"op_mult"`
	OpDivTime      int `json:"op_div"`
	OpAgentTimeOut int `json:"op_agent_timeout"`
}

type Task struct {
	Expression Expression
	Settings   Settinsg
}
