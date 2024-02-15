package handlers

import (
	"net/http"
	"strconv"

	"github.com/kozlovsv/evaluator/server/pkg/models"
)

type SetResultHandler struct {
	store      models.ExpressionStore
	storeAgent models.AgentStore
}

func NewSetResultHandlerr(s models.ExpressionStore, storeAgent models.AgentStore) *SetResultHandler {
	return &SetResultHandler{
		store:      s,
		storeAgent: storeAgent,
	}
}

func (h *SetResultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 [Set Result] Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "400 [SetResult] BadRequest. ID not integer."+r.URL.Query().Get("id"), http.StatusBadRequest)
		return
	}

	result := r.URL.Query().Get("res")

	err = h.store.SetExpressionResult(id, result)

	if err != nil {
		http.Error(w, "500 [SetResult] InternalServerError. Error save result to storage"+err.Error(), http.StatusInternalServerError)
		return
	}

	agentName := r.URL.Query().Get("agent")
	h.storeAgent.Add(agentName)
	h.storeAgent.SetCurrentOp(agentName, "")

	w.WriteHeader(http.StatusOK)
}
