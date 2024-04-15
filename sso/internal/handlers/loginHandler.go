package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kozlovsv/evaluator/sso/internal/lib/logger/sl"
	models "github.com/kozlovsv/evaluator/sso/internal/models"
	auth "github.com/kozlovsv/evaluator/sso/internal/services"
)

type LoginHandler struct {
	auth auth.AuthIntfs
	log  *slog.Logger
}

func NewLoginHandler(auth auth.AuthIntfs, log *slog.Logger) *LoginHandler {
	return &LoginHandler{
		auth: auth,
		log:  log,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lg := h.log.With(slog.String("op", "handlers.loginHundler"))

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

	token, err := h.auth.Login(r.Context(), regstr.Email, regstr.Password)

	if err != nil {
		// Ошибку auth.ErrInvalidCredentials мы создадим ниже
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "401 Bad auth", http.StatusUnauthorized)
			lg.Info("Неверный логин проль", sl.Err(err))
			return
		}

		http.Error(w, "500 Internal Error", http.StatusInternalServerError)
		lg.Info("Failed login user", sl.Err(err))
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
