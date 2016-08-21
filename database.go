package sorm

import (
	"database/sql"
	"fmt"
	"time"
)

type database struct {
	dbtype string  // just now only support "mysql"
	dsn    string  // connection string
	db     *sql.DB // underlying sql connection
}

func (db *database) open(conn string) (err error) {
	if conn == "" {
		return fmt.Errorf("invalid db connection string")
	}
	db.dsn = conn

	db.db, err = sql.Open("mysql", conn)
	if err != nil {
		return err
	}

	err = db.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *database) Exec(sql string, args ...interface{}) (res sql.Result, err error) {
	if db.db != nil {
		return db.db.Exec(sql, args...)
	}
	return nil, fmt.Errorf("db is not opened")
}

func (db *database) Close() (err error) {
	if db.db != nil {
		err = db.db.Close()
		if err == nil {
			db.db = nil
		}
	}
	return err
}

func (db *database) BindTable(tn string) (t Table, err error) {
	tbl := &table{db: db, name: tn}
	err = tbl.query("select * from " + tn)
	if err != nil {
		return nil, err
	}
	return tbl, nil
}

func (db *database) CreateQuery(sql string) (q Query, err error) {
	if db.db != nil {
		qr := &query{sql: sql}
		qr.stmt, err = db.db.Prepare(sql)
		if err != nil {
			return nil, err
		}
		q = qr
	}
	return q, err
}

func (db *database) SetConnMaxLifetime(d time.Duration) {
	if db.db != nil {
		db.db.SetConnMaxLifetime(d)
	}
}

func (db *database) SetMaxIdleConns(n int) {
	if db.db != nil {
		db.db.SetMaxIdleConns(n)
	}
}

func (db *database) SetMaxOpenConns(n int) {
	if db.db != nil {
		db.db.SetMaxOpenConns(n)
	}
}
