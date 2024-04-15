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
	UserId         int64
}
