package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kozlovsv/evaluator/sso/internal/lib/logger/sl"
	models "github.com/kozlovsv/evaluator/sso/internal/models"
	auth "github.com/kozlovsv/evaluator/sso/internal/services"
	storage "github.com/kozlovsv/evaluator/sso/internal/storage"
)

type RegisterHandler struct {
	auth auth.AuthIntfs
	log  *slog.Logger
}

func NewRegisterHandler(auth auth.AuthIntfs, log *slog.Logger) *RegisterHandler {
	return &RegisterHandler{
		auth: auth,
		log:  log,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.registerHundler"

	lg := h.log.With(slog.String("op", op))

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Recipe object that will be populated from json payload
	var regstr models.Register

	if err := json.NewDecoder(r.Body).Decode(&regstr); err != nil {
		http.Error(w, "400 Error decode JSON"+err.Error(), http.StatusBadRequest)
		lg.Info("Error decode JSON", sl.Err(err))
		return
	}

	if regstr.Email == "" {
		http.Error(w, "400 BadRequest Email Req", http.StatusBadRequest)
		return
	}

	if regstr.Password == "" {
		http.Error(w, "405 BadRequest Password Req", http.StatusBadRequest)
		return
	}

	_, err := h.auth.RegisterNewUser(r.Context(), regstr.Email, regstr.Password)

	if err != nil {
		// Ошибку storage.ErrUserExists мы создадим ниже
		if errors.Is(err, storage.ErrUserExists) {
			http.Error(w, "405 user already exists", http.StatusBadRequest)
			lg.Info("User already exists", sl.Err(err))
			return
		}

		http.Error(w, "500 failed to register user", http.StatusInternalServerError)
		lg.Info("Failed to register user", sl.Err(err))
		return
	}

	//После добавления юзера сразу авторизуемся
	token, err := h.auth.Login(r.Context(), regstr.Email, regstr.Password)

	if err != nil {
		http.Error(w, "500 Internal Error", http.StatusInternalServerError)
		lg.Info("Failed login after register user", sl.Err(err))
		return
	}

	tokenModel := models.Token{
		Token:     token,
		UserEmail: regstr.Email,
	}

	jsonBytes, _ := json.Marshal(tokenModel)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
