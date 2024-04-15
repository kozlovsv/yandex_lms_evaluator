package models

type Token struct {
	Token     string `json:"access_token"`
	UserEmail string `json:"user_email"`
}
