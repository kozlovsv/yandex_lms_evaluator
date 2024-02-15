package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kozlovsv/evaluator/server/pkg/models"
)

type GetNewExpressionHandler struct {
	storeExpression models.ExpressionStore
	storeSettings   models.SettingsStore
	storeAgent      models.AgentStore
}

func NewGetNewExpressionHandler(storeExpression models.ExpressionStore, storeSettings models.SettingsStore, storeAgent models.AgentStore) *GetNewExpressionHandler {
	return &GetNewExpressionHandler{
		storeExpression: storeExpression,
		storeSettings:   storeSettings,
		storeAgent:      storeAgent,
	}
}

func (h *GetNewExpressionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	agentName := r.URL.Query().Get("agent")
	h.storeAgent.Add(agentName)

	exp, err := h.storeExpression.GetNewExpression()

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			h.storeAgent.SetAvailable(agentName)
			return
		}
		log.Println("[ERROR]", "[Get New Expression]", "Error load expression from storage", err.Error())
		http.Error(w, "500  Internal Server Error", http.StatusInternalServerError)
		return
	}

	settings, err := h.storeSettings.Get()

	if err != nil {
		log.Println("[ERROR]", "[Get New Expression]", "Error load settings from storage", err.Error())
		http.Error(w, "500  Internal Server Error", http.StatusInternalServerError)
		return
	}

	task := models.Task{
		Expression: exp,
		Settings:   settings,
	}
	jsonBytes, err := json.Marshal(task)
	if err != nil {
		log.Println("[ERROR]", "[Get New Expression]", "Error code to JSON", err.Error())
		http.Error(w, "500  Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.storeAgent.SetCurrentOp(agentName, exp.Value)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
