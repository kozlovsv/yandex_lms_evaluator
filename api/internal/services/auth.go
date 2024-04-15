package auth

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/kozlovsv/evaluator/api/config"
	"github.com/kozlovsv/evaluator/api/internal/lib/logger/sl"
)

// Тот самый интерфейс авторизации
type AuthIntfs interface {
	ParseToken(tokenStr string) (userId int64, err error)
}

type Auth struct {
	log       *slog.Logger
	jwtConfig config.JWTConfig
}

func New(
	log *slog.Logger,
	jwtConfig config.JWTConfig,
) *Auth {
	return &Auth{
		log:       log,
		jwtConfig: jwtConfig,
	}
}

var (
	ErrBadJWTToken  = errors.New("Bad token format")
	InvalidJWTToken = errors.New("Invalid token")
)

type CustomClaims struct {
	UserID int64  `json:"uid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewToken creates new JWT token for given user and app.
func (a *Auth) ParseToken(tokenStr string) (int64, error) {

	log := a.log.With(
		slog.String("op", "Auth.ParseToken"),
		slog.String("token", tokenStr),
	)

	// Парсинг и проверка токена
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(a.jwtConfig.Secret), nil },
	)

	if err != nil {
		log.Warn("Bad token format", sl.Err(err))
		return 0, ErrBadJWTToken
	}

	// Проверка валидности токена
	if !token.Valid {
		log.Warn("Token invalid", slog.Any("token", token))
		return 0, InvalidJWTToken
	}

	claims, ok := token.Claims.(*CustomClaims)

	if !ok {
		log.Warn("Error loading token claims", slog.Any("token", token))
		return 0, InvalidJWTToken
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
		log.Warn("Token Expired", slog.Any("ExpiresAt", claims.ExpiresAt))
		return 0, InvalidJWTToken
	}

	return claims.UserID, nil
}
