package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kozlovsv/evaluator/api/internal/storage"
)

type AgentHandler struct {
	store *storage.AgentStore
}

func NewAgentHandler(s *storage.AgentStore) *AgentHandler {
	return &AgentHandler{
		store: s,
	}
}

func (h *AgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		h.GetAgents(w, r)
		return
	default:
		log.Println("[ERROR]", "AgentHandler", "405 Method Not Allowed")
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *AgentHandler) GetAgents(w http.ResponseWriter, r *http.Request) {
	list, err := h.store.List()
	if err != nil {
		log.Println("[ERROR]", "[Agent Handler]", "Ошибка получения списка из базы", err.Error())
		http.Error(w, "500 Internal Error", http.StatusInternalServerError)
	}
	jsonBytes, err := json.Marshal(list)
	if err != nil {
		log.Println("[ERROR]", "[Agent Handler]", "Error decode JSON", err.Error())
		http.Error(w, "500  Error decode JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
