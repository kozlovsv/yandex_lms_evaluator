package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kozlovsv/evaluator/server/pkg/models"
)

type SetErrorHandler struct {
	store      models.ExpressionStore
	storeAgent models.AgentStore
}

func NewSetErrorHandlerr(s models.ExpressionStore, storeAgent models.AgentStore) *SetErrorHandler {
	return &SetErrorHandler{
		store:      s,
		storeAgent: storeAgent,
	}
}

func (h *SetErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("[ERROR] 405 [Set Error] Method Not Allowed")
		http.Error(w, "405 [Set Error] Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("[ERROR] 500 [SetError] BadRequest. ID not integer." + r.URL.Query().Get("id"))
		http.Error(w, "250 [SetError] BadRequest. ID not integer."+r.URL.Query().Get("id"), http.StatusBadGateway)
		return
	}

	errorStr := r.URL.Query().Get("err")

	err = h.store.SetExpressionError(id, errorStr)

	if err != nil {
		log.Println("[ERROR] 500 [SetResult] InternalServerError. Error save Error to storage" + err.Error())
		http.Error(w, "500 [SetResult] InternalServerError. Error save Error to storage"+err.Error(), http.StatusInternalServerError)
		return
	}

	agentName := r.URL.Query().Get("agent")
	h.storeAgent.Add(agentName)
	h.storeAgent.SetCurrentOp(agentName, "")

	w.WriteHeader(http.StatusOK)
}
