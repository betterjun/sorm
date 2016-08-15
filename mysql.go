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

func (db *database) Query(sql string, objs interface{}, args ...interface{}) (err error) {
	if db.db != nil {
		rows, err := db.db.Query(sql, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		var cols []string
		cols, err = rows.Columns()
		if err != nil {
			return err
		}

		val := reflect.ValueOf(objs)
		sInd := reflect.Indirect(val)

		if val.Kind() != reflect.Ptr || sInd.Kind() != reflect.Slice {
			return fmt.Errorf("<Database.Query> output arg must be use ptr slice")
		}

		var scanArgs []interface{}
		var ind reflect.Value
		etyp := sInd.Type().Elem()
		if etyp.Kind() == reflect.Struct {
			ind = reflect.New(sInd.Type().Elem()).Elem()
			scanArgs = getScanFieldFromStruct(ind, cols)
			if scanArgs == nil {
				return fmt.Errorf("no fields found in the objptr")
			}
		} else {
			if len(cols) > 1 {
				return fmt.Errorf("the query returns multi coloums, please passing in a struct slice")
			}

			ind = reflect.New(etyp).Elem()
			scanArgs = []interface{}{ind.Addr().Interface()}
		}

		sIndCopy := sInd
		err = fmt.Errorf("no records found")
		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				break
			}

			sIndCopy = reflect.Append(sIndCopy, ind)
		}

		sInd.Set(sIndCopy)
		return err
	}
	return fmt.Errorf("db is not opened")
}

func (db *database) Exec(sql string, args ...interface{}) (res sql.Result, err error) {
	if db.db != nil {
		return db.db.Exec(sql, args...)
	}
	return res, fmt.Errorf("db is not opened")
}

type query struct {
	sql   string
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
	
	if q.rows.Next() {
		scanArgs := getScanFields(obj, q.cols)
		if scanArgs == nil {
			return fmt.Errorf("no fields found in the objptr")
		}

		return q.rows.Scan(scanArgs...)
	} else {
		q.rows.Close()
		return fmt.Errorf("end of query results")
	}
}

func (q *query) All(objs interface{}) (err error) {
	val := reflect.ValueOf(objs)
	sInd := reflect.Indirect(val)
	if val.Kind() != reflect.Ptr || sInd.Kind() != reflect.Slice {
		return fmt.Errorf("<Query.All> output arg must be use ptr slice")
	}

	if q.rows == nil {
		return fmt.Errorf("query is not executed")
	}
	cols, err := q.rows.Columns()
	if err != nil {
		return err
	}
	defer q.rows.Close()
	
	var scanArgs []interface{}
	var ind reflect.Value
	etyp := sInd.Type().Elem()
	if etyp.Kind() == reflect.Struct {
		ind = reflect.New(sInd.Type().Elem()).Elem()
		scanArgs = getScanFieldFromStruct(ind, cols)
		if scanArgs == nil {
			return fmt.Errorf("no fields found in the objptr")
		}
	} else {
		if len(cols) > 1 {
			return fmt.Errorf("the query returns multi coloums, please passing in a struct slice")
		}

		ind = reflect.New(etyp).Elem()
		scanArgs = []interface{}{ind.Addr().Interface()}
	}

	sIndCopy := sInd
	err = fmt.Errorf("no records found")
	for q.rows.Next() {
		err = q.rows.Scan(scanArgs...)
		if err != nil {
			break
		}

		sIndCopy = reflect.Append(sIndCopy, ind)
	}

	// ret may take back some records, even though there is an error.
	sInd.Set(sIndCopy)
	return err
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
