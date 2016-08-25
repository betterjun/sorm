package sorm

import (
	"database/sql"
	"fmt"
)

type query struct {
	sql  string
	stmt *sql.Stmt
}

func (q *query) Exec(args ...interface{}) (res Result, err error) {
	if q.stmt == nil {
		return nil, fmt.Errorf("query is not initialized")
	}

	if printSql {
		fmt.Printf("Query.Exec: %v, args %v\n", q.sql, args)
	}
	rows, err := q.stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	res = &result{rows: rows}
	return res, nil
}

func (q *query) Close() (err error) {
	if q.stmt != nil {
		err = q.stmt.Close()
		q.stmt = nil
	}
	return err
}
