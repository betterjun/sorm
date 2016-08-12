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
	t = &table{db: db.db, name: tn}
	return t
}

func (db *database) CreateQuery(sql string, model interface{}) (q Query, err error) {
	if db.db != nil {
		qr := &query{model: model, sql: sql}
		qr.stmt, err = db.db.Prepare(sql)
		if err != nil {
			return nil, err
		}
		q = qr
	}
	return q, err
}

func (db *database) QueryRow(sql string, objptr interface{}, args ...interface{}) (err error) {
	if db.db != nil {
		rows, err := db.db.Query(sql, args...)
		if err != nil {
			return err
		}
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
	sql   string
	model interface{}
	stmt  *sql.Stmt
	rows  *sql.Rows
	cols  []string
}

func (q *query) ExecuteQuery(args ...interface{}) (err error) {
	if q.rows != nil {
		q.rows.Close()
	}
	q.rows, err = q.stmt.Query(args...)
	return err
}

func (q *query) Next(obj interface{}) (err error) {
	if q.rows == nil {
		return fmt.Errorf("query is not executed")
	}
	if q.cols == nil {
		q.cols, err = q.rows.Columns()
		if err != nil {
			return err
		}
	}

	scanArgs := getScanFields(obj, q.cols)
	if scanArgs == nil {
		return fmt.Errorf("no valid model is created")
	}
	if q.rows.Next() {
		err = q.rows.Scan(scanArgs...)
		if err != nil {
			q.rows.Close()
		}
	} else {
		err = fmt.Errorf("end of query results")
		q.rows.Close()
	}
	return err
}

func (q *query) All() (ret []interface{}, err error) {
	if q.rows == nil {
		return nil, fmt.Errorf("query is not executed")
	}
	if q.cols == nil {
		q.cols, err = q.rows.Columns()
		if err != nil {
			return nil, err
		}
	}

	scanArgs := getScanFields(q.model, q.cols)
	if scanArgs == nil {
		return nil, fmt.Errorf("no valid model is created")
	}
	defer q.rows.Close()
	// if not in the below loop, err should be the eof.
	err = fmt.Errorf("end of query results")
	for q.rows.Next() {
		err = q.rows.Scan(scanArgs...)
		if err != nil {
			break
		}
		ret = append(ret, reflect.Indirect(reflect.ValueOf(q.model)))
	}
	// ret may take back some records, even though there is an error.
	return ret, err
}

type table struct {
	query
	name string
	db   *sql.DB
}

func (t *table) ExecuteQuery(args ...interface{}) (err error) {
	if t.db != nil {
		t.sql = fmt.Sprintf("select * from %v", t.name)
		t.stmt, err = t.db.Prepare(t.sql)
		if err != nil {
			return err
		}

		if t.rows != nil {
			t.rows.Close()
		}
		t.rows, err = t.stmt.Query(args...)

		// todo: model is a map, shall pass a struct when bind the table?
	}

	return err
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
