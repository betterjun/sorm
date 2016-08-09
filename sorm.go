package sorm

import (
	"database/sql"
	"fmt"
)

//
func NewDatabase(dbtype, conn string) (db Database) {
	switch dbtype {
	case "mysql":
		d := &database{}
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
	BindTable(tn string) Table
	CreateQuery(sql string) (Query, error)

	QueryRow(sql string, obj interface{}, args ...interface{}) error
	Query(sql string, objs []interface{}, args ...interface{}) error

	// if args[0] is a struct, use it's field
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Close() error
	//SetConnMaxLifetime(d time.Duration)
	//SetMaxIdleConns(n int)
	//SetMaxOpenConns(n int)
}

type Table interface {
	Query
	Insert(obj interface{}) (sql.Result, error)
	// by pk
	Delete(obj interface{}) (sql.Result, error)
	Update(obj interface{}) (sql.Result, error)
	//Count(obj interface{}) (int64, error)

	//Drop() error
	//Truncate() error
	/*
		SetPrimaryKey(name string)
		SetField(fn string, ignore bool)
		SetIgnoredField(fn string)
	*/
	// supported tags
	// `fn:"field_name" pk:"1"`
}

type Query interface {
	Exec(args ...interface{}) (sql.Result, error)
	Next(obj interface{}) error
	QueryRow(obj interface{}) error
	Query(objs []interface{}) error

	//Columns() ([]string, error)
}
