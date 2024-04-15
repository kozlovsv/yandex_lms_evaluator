package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kozlovsv/evaluator/api/internal/models"
	"github.com/kozlovsv/evaluator/api/internal/storage"
)

// ExpressionHandler implements http.Handler and dispatch request to the store
type SettingsHandler struct {
	store *storage.SettingsStore
}

func NewSettingsHandler(s *storage.SettingsStore) *SettingsHandler {
	return &SettingsHandler{
		store: s,
	}
}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		h.SetSettings(w, r)
		return
	case r.Method == http.MethodGet:
		h.GetSettings(w, r)
		return
	default:
		log.Println("[ERROR]", "SettingsHandler", "405 Method Not Allowed")
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *SettingsHandler) SetSettings(w http.ResponseWriter, r *http.Request) {
	// Recipe object that will be populated from json payload
	var settings models.Settinsg

	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		log.Println("[ERROR]", "[SettingsHandler]", "Error decode JSON", err.Error())
		http.Error(w, "400 Bad request", http.StatusBadRequest)
		return
	}

	if err := h.store.Set(settings); err != nil {
		log.Println("[ERROR]", "[SettingsHandler]", "Save settings to storage", err.Error())
		http.Error(w, "500  Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.store.Get()
	if err != nil {
		log.Println("[ERROR]", "[SettingsHandler]", "Error load row from storage", err.Error())
		http.Error(w, "400 Bad request", http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(settings)
	if err != nil {
		log.Println("[ERROR]", "[SettingsHandler]", "Error code JSON from sending", err.Error())
		http.Error(w, "500  Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
