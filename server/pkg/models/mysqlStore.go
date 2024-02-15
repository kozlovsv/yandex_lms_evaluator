package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("mysql", "root:testerum@tcp(db:3306)/evaluator?parseTime=true")
}
