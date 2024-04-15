package models

import "time"

type Agent struct {
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	LastPing  time.Time `json:"last_ping"`
	CurrentOp string    `json:"current_op"`
}
