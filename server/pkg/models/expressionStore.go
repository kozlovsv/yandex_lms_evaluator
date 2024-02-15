package models

import (
	"database/sql"
	"fmt"
	"log"
)

type ExpressionStore interface {
	Add(expression Expression) error
	Get(id int) (Expression, error)
	List() ([]Expression, error)
	GetNewExpression() (Expression, error)
	SetExpressionResult(id int, result string) error
	SetExpressionStatus(id int, status int) error
	SetExpressionError(id int, error string) error
	UpFrozenExpressions(timeout int) error
}

type RowScanner interface {
	Scan(src ...any) error
}

var expressionFieldsSet = "`id`, `value`, `status`, `result`, `idempotency_key`, `updated_at`, `created_at`"

type ExpressionMySqlStore struct {
	DB *sql.DB
}

func NewExpressionStore(db *sql.DB) *ExpressionMySqlStore {
	return &ExpressionMySqlStore{
		DB: db,
	}
}

func (s ExpressionMySqlStore) Add(expression Expression) error {
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

	insert, err := s.DB.Query("INSERT INTO expression (value, idempotency_key) VALUES (?, ?)", expression.Value, expression.IdempotencyKey)

	if err != nil {
		return err
	}

	defer insert.Close()
	return nil
}

func (s ExpressionMySqlStore) getExpression(row RowScanner) (Expression, error) {
	expression := Expression{}
	err := row.Scan(&expression.Id, &expression.Value, &expression.Status, &expression.Result, &expression.IdempotencyKey, &expression.UpdatedAt, &expression.CreatesAt)
	return expression, err
}

func (s ExpressionMySqlStore) Get(id int) (Expression, error) {
	row := s.DB.QueryRow(fmt.Sprintf("SELECT %s FROM expression WHERE id = ?", expressionFieldsSet), id)
	return s.getExpression(row)
}

func (s ExpressionMySqlStore) List() ([]Expression, error) {
	var res []Expression
	query := fmt.Sprintf("SELECT %s FROM expression ORDER BY id DESC", expressionFieldsSet)
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

func (s ExpressionMySqlStore) SetExpressionResult(id int, result string) error {
	_, err := s.DB.Exec("UPDATE expression SET result = ?, status = 3, updated_at = CURRENT_TIMESTAMP() WHERE id = ? AND status <> 3", result, id)
	return err

}

func (s ExpressionMySqlStore) SetExpressionError(id int, error string) error {
	_, err := s.DB.Exec("UPDATE expression SET result_text = ?, status = 2, updated_at = CURRENT_TIMESTAMP() WHERE id = ? AND status <> 3", error, id)
	return err
}

func (s ExpressionMySqlStore) SetExpressionStatus(id int, status int) error {
	_, err := s.DB.Exec("UPDATE expression SET status = ?, updated_at = CURRENT_TIMESTAMP() WHERE id = ?", status, id)
	return err
}

func (s ExpressionMySqlStore) GetNewExpression() (Expression, error) {
	tx, _ := s.DB.Begin()
	row := tx.QueryRow(fmt.Sprintf("SELECT %s FROM expression WHERE status = 0 ORDER BY id", expressionFieldsSet))
	exp, err := s.getExpression(row)
	if err != nil {
		tx.Rollback()
		return exp, err
	}

	err = s.SetExpressionStatus(exp.Id, 1)

	if err != nil {
		tx.Rollback()
		return Expression{}, err
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return Expression{}, err
	}

	return exp, nil
}

func (s ExpressionMySqlStore) UpFrozenExpressions(timeout int) error {
	stm, err := s.DB.Exec("UPDATE expression SET status = 0, updated_at = CURRENT_TIMESTAMP() WHERE status = 1 AND updated_at <= DATE_SUB(NOW(), INTERVAL ? * 1000 MICROSECOND)", timeout)
	if err == nil {
		if cnt, _ := stm.RowsAffected(); cnt > 0 {
			log.Println("[INFO] Разморожено задач:", cnt)
		}
	}
	return nil
}
