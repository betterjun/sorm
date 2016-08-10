package sorm

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type database struct {
	basedb
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

func (db *database) Close() (err error) {
	if db.db != nil {
		err = db.db.Close()
		if err == nil {
			db.db = nil
		}
	}
	return err
}

func (db *database) BindTable(tn string) (t Table) {
	t = &table{}
	return t
}

func (db *database) CreateQuery(sql string) (q Query, err error) {
	q = &query{}
	return q, err
}

func (db *database) QueryRow(sql string, objptr interface{}, args ...interface{}) (err error) {
	if db.db != nil {
		rows, err := db.db.Query(sql, args...)
		defer rows.Close()
		if rows.Next() {
			var cols []string
			cols, err = rows.Columns()
			if err != nil {
				return err
			}
			scanArgs := getScanFields(objptr, cols)
			if scanArgs == nil {
				return fmt.Errorf("no fields found in the objptr")
			}

			err = rows.Scan(scanArgs...)
			//fmt.Println("objptr=", objptr)
			return err
		} else {
			return fmt.Errorf("no record found")
		}
	}
	return fmt.Errorf("db is not opened")
}

func (db *database) Query(sql string, model interface{}, args ...interface{}) (objs []interface{}, err error) {
	if db.db != nil {
		rows, err := db.db.Query(sql, args...)
		defer rows.Close()

		var cols []string
		cols, err = rows.Columns()
		if err != nil {
			return nil, err
		}

		scanArgs := getScanFields(model, cols)
		if scanArgs == nil {
			return nil, fmt.Errorf("no fields found in the objptr")
		}

		for _, v := range scanArgs {
			fmt.Println("sorm :", v)
		}

		ret := make([]interface{}, 0)
		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				break
			}
			ret = append(ret, reflect.Indirect(reflect.ValueOf(model)))
		}
		return ret, err
	}
	return nil, fmt.Errorf("db is not opened")
}

func (db *database) Exec(sql string, args ...interface{}) (res sql.Result, err error) {
	if db.db != nil {
		return db.db.Exec(sql, args...)
	}
	return res, fmt.Errorf("db is not opened")
}

type query struct {
	dbtype string // just now only support "mysql"
	dsn    string // connection string
	db     *sql.DB
}

func (q *query) Exec(args ...interface{}) (res sql.Result, err error) {
	return
}

func (q *query) Next(obj interface{}) (err error) {
	return
}

func (q *query) QueryRow(obj interface{}) (err error) {
	return
}

func (q *query) Query(model interface{}) (objs []interface{}, err error) {
	return
}

type table struct {
	dbtype string // just now only support "mysql"
	dsn    string // connection string
	db     *sql.DB
}

func (t *table) Exec(args ...interface{}) (res sql.Result, err error) {
	return
}

func (t *table) Next(obj interface{}) (err error) {
	return
}

func (t *table) QueryRow(obj interface{}) (err error) {
	return
}

func (t *table) Query(model interface{}) (objs []interface{}, err error) {
	return
}

func (t *table) Insert(obj interface{}) (res sql.Result, err error) {
	return
}

func (t *table) Delete(obj interface{}) (res sql.Result, err error) {
	return
}

func (t *table) Update(obj interface{}) (res sql.Result, err error) {
	return
}
