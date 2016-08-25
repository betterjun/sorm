package sorm

import (
	"database/sql"
	"fmt"
	"time"
)

var printSql bool = false

func PrintSql(yes bool) {
	printSql = yes
}

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
}

type Query interface {
	// need first call Exec
	Exec(args ...interface{}) (res Result, err error)
	Close() error
}

type Result interface {
	Next(obj interface{}, args ...interface{}) error
	All(objs interface{}) error
	ColumnNames() ([]string, error)
	Close() error

	// ColumnByIndex(col int) ([]interface{}, error)
	// ColumnByName(col string) ([]interface{}, error)
}
