package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	config "github.com/kozlovsv/evaluator/sso/config"
	"github.com/kozlovsv/evaluator/sso/internal/lib/jwt"
	"github.com/kozlovsv/evaluator/sso/internal/lib/logger/sl"
	"github.com/kozlovsv/evaluator/sso/internal/models"
	"github.com/kozlovsv/evaluator/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// Тот самый интерфейс авторизации
type AuthIntfs interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type Auth struct {
	log       *slog.Logger
	storage   UserStorage
	jwtConfig config.JWTConfig
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserStorage interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
	User(ctx context.Context, email string) (models.User, error)
}

func New(
	log *slog.Logger,
	storage UserStorage,
	jwtConfig config.JWTConfig,
) *Auth {
	return &Auth{
		storage:   storage,
		log:       log,
		jwtConfig: jwtConfig,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.storage.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, a.jwtConfig.TokenTTL, a.jwtConfig.Secret)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.storage.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
