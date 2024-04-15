package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/kozlovsv/evaluator/api/internal/models"
	"github.com/kozlovsv/evaluator/api/internal/storage"
)

var (
	ExpressionUrlRe       = regexp.MustCompile(`^/expressions/*$`)
	ExpressionUrlReWithID = regexp.MustCompile(`^/expressions/(\d*)$`)
)

// ExpressionHandler implements http.Handler and dispatch request to the store
type ExpressionHandler struct {
	store *storage.ExpressionStore
}

func NewExpressionHandler(s *storage.ExpressionStore) *ExpressionHandler {
	return &ExpressionHandler{
		store: s,
	}
}

func (h *ExpressionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		h.CreateExpression(w, r)
		return
	case r.Method == http.MethodGet && ExpressionUrlRe.MatchString(r.URL.Path):
		h.ListExpressions(w, r)
		return
	case r.Method == http.MethodGet && ExpressionUrlReWithID.MatchString(r.URL.Path):
		h.GetExpression(w, r)
		return
	default:
		http.Error(w, "405 [New Expression] Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *ExpressionHandler) CreateExpression(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextUserIDKey).(int64)

	// Recipe object that will be populated from json payload
	var exp models.Expression

	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		http.Error(w, "400 [CreateExpression] Error decode JSON "+err.Error(), http.StatusBadRequest)
		return
	}
	exp.UserId = userId
	if err := h.store.Add(exp); err != nil {
		http.Error(w, "500 [CreateExpression] "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ExpressionHandler) ListExpressions(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ContextUserIDKey).(int64)
	list, err := h.store.List(userId)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		http.Error(w, "500 [ListExpressions] Load From DB List"+err.Error(), http.StatusInternalServerError)
	}
	jsonBytes, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "500 [ListExpressions] Error decode JSON"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *ExpressionHandler) GetExpression(w http.ResponseWriter, r *http.Request) {
	matches := ExpressionUrlReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		log.Println("[ERROR]", "[GetExpression]", "Bad URL format. Not Found ID")
		http.Error(w, "400 [GetExpression] BadRequest. Bad URL format. Not Found ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil || id < 1 {
		log.Println("[ERROR]", "[GetExpression]", "Bad URL format. ID NOT Integer OR < -1", matches[1])
		http.Error(w, "400 [GetExpression] BadRequest. Bad URL format. ID NOT Integer OR < -1. "+matches[1], http.StatusBadRequest)
		return
	}

	exp, err := h.store.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Println("[ERROR]", "[GetExpression]", "Error load row from storage", err.Error())
		http.Error(w, "400 [GetExpression] Error load row from storage", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value(ContextUserIDKey).(int64)
	if userId != exp.UserId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	jsonBytes, err := json.Marshal(exp)
	if err != nil {
		log.Println("[ERROR]", "[GetExpression]", "Error code JSON from sending", err.Error())
		http.Error(w, "500 [GetExpression] Error decode JSON"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
