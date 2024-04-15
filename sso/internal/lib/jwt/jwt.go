package jwt

import (
	"time"

	"github.com/kozlovsv/evaluator/sso/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	Uid   int64  `json:"uid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// NewToken creates new JWT token for given user and app.
func NewToken(user models.User, duration time.Duration, secret string) (string, error) {
	claims := customClaims{
		Uid:   user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
