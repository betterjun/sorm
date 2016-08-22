package sorm

import (
	"database/sql"
	"fmt"
	"time"
)

//
func NewDatabase(dbtype, conn string) (db Database) {
	switch dbtype {
	case "mysql":
		d := &database{}
		d.dbtype = dbtype
		err := d.open(conn)
		if err != nil {
			fmt.Printf("sorm create database connection failed:%v\n", err)
			return nil
		}
		return d
	default:
		fmt.Printf("sorm create database connection failed: unsurpported dbtype %v\n", dbtype)
		return nil
	}
}

type Database interface {
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Close() error

	BindTable(tn string) (Table, error)
	CreateQuery(sql string) (Query, error)

	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
}

type Table interface {
	/*
		Result
		Insert(obj interface{}) (sql.Result, error)
		// by pk
		Delete(obj interface{}) (sql.Result, error)
		Update(obj interface{}) (sql.Result, error)
	*/

	// Refactor the methods as below?
	// type Filter map[string]interface{}
	//Query(obj, filter)
	//Delete(obj, filter)
	//Update(obj, filter)

	// will insert by the column order
	Insert(values ...interface{}) (sql.Result, error)
	// will only insert the columns the table has
	//Insert(value map[string]interface{})
	//Insert(value struct)

	Delete(filter string) (sql.Result, error)

	Update(filter string, value interface{}) (sql.Result, error)
	//Update(filter string, value map[string]interface{})
	//Update(filter string, value struct)

	// will select all columns
	Query(filter string) (Result, error)

	//Drop() error
	//Truncate() error
}

type Query interface {
	// need first call Exec
	Exec(args ...interface{}) (res Result, err error)
	Close() error
}

/*
The design gives a good reuse of Query, can have multi result sets by each Exec call.
res := Query.Exec()
res.Next()
res.All()

res.Filter()
*/
type Result interface {
	Next(obj interface{}, args ...interface{}) error
	All(objs interface{}) error
	Close() error

	//Count(filter interface{}) (int64, error)
	//Columns() ([]string, error)
	//Filter()
}
