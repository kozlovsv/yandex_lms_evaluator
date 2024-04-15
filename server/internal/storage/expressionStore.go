package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kozlovsv/evaluator/server/internal/models"
)

type RowScanner interface {
	Scan(src ...any) error
}

var expressionFieldsSet = "`id`, `value`"

type ExpressionStore struct {
	DB *sql.DB
}

func NewExpressionStore(db *sql.DB) *ExpressionStore {
	return &ExpressionStore{
		DB: db,
	}
}

func (s ExpressionStore) GetNewExpression() (models.Expression, error) {
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
		return models.Expression{}, err
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return models.Expression{}, err
	}

	return exp, nil
}

func (s ExpressionStore) getExpression(row RowScanner) (models.Expression, error) {
	expression := models.Expression{}
	err := row.Scan(&expression.Id, &expression.Value)
	return expression, err
}

func (s ExpressionStore) SetExpressionResult(id int, result string) error {
	_, err := s.DB.Exec("UPDATE expression SET result = ?, status = 3, updated_at = CURRENT_TIMESTAMP() WHERE id = ? AND status <> 3", result, id)
	return err

}

func (s ExpressionStore) SetExpressionError(id int, error string) error {
	_, err := s.DB.Exec("UPDATE expression SET result_text = ?, status = 2, updated_at = CURRENT_TIMESTAMP() WHERE id = ? AND status <> 3", error, id)
	return err
}

func (s ExpressionStore) SetExpressionStatus(id int, status int) error {
	_, err := s.DB.Exec("UPDATE expression SET status = ?, updated_at = CURRENT_TIMESTAMP() WHERE id = ?", status, id)
	return err
}

func (s ExpressionStore) UpFrozenExpressions(timeout int) error {
	stm, err := s.DB.Exec("UPDATE expression SET status = 0, updated_at = CURRENT_TIMESTAMP() WHERE status = 1 AND updated_at <= DATE_SUB(NOW(), INTERVAL ? * 1000 MICROSECOND)", timeout)
	if err == nil {
		if cnt, _ := stm.RowsAffected(); cnt > 0 {
			log.Println("[INFO] Разморожено задач:", cnt)
		}
	}
	return nil
}
