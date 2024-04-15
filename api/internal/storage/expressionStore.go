package storage

import (
	"database/sql"
	"fmt"

	"github.com/kozlovsv/evaluator/api/internal/models"
)

type RowScanner interface {
	Scan(src ...any) error
}

var expressionFieldsSet = "`id`, `value`, `status`, `result`, `idempotency_key`, `updated_at`, `created_at`, `user_id`"

type ExpressionStore struct {
	DB *sql.DB
}

func NewExpressionStore(db *sql.DB) *ExpressionStore {
	return &ExpressionStore{
		DB: db,
	}
}

func (s ExpressionStore) Add(expression models.Expression) error {
	//Проверяем может уже есть такой запрос
	id := 0
	err := s.DB.QueryRow("SELECT id FROM expression WHERE idempotency_key = ? LIMIT 1", expression.IdempotencyKey).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	//Если есть такой ID то ничего добавлять не нужно
	if id > 0 {
		return nil
	}

	insert, err := s.DB.Query("INSERT INTO expression (value, idempotency_key, user_id) VALUES (?, ?, ?)", expression.Value, expression.IdempotencyKey, expression.UserId)

	if err != nil {
		return err
	}

	defer insert.Close()
	return nil
}

func (s ExpressionStore) getExpression(row RowScanner) (models.Expression, error) {
	expression := models.Expression{}
	err := row.Scan(&expression.Id, &expression.Value, &expression.Status, &expression.Result, &expression.IdempotencyKey, &expression.UpdatedAt, &expression.CreatesAt, &expression.UserId)
	return expression, err
}

func (s ExpressionStore) Get(id int) (models.Expression, error) {
	row := s.DB.QueryRow(fmt.Sprintf("SELECT %s FROM expression WHERE id = ?", expressionFieldsSet), id)
	return s.getExpression(row)
}

func (s ExpressionStore) List(userId int64) ([]models.Expression, error) {
	var res []models.Expression
	query := fmt.Sprintf("SELECT %s FROM expression WHERE user_id = %d ORDER BY id DESC", expressionFieldsSet, userId)
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item, err := s.getExpression(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}
