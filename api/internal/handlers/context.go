package handlers

type ContextKey string

const (
	ContextJWTKey    ContextKey = "JWTToken"
	ContextUserIDKey ContextKey = "UserID"
)
