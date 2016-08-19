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

type basedb struct {
	dbtype string // just now only support "mysql"
	dsn    string // connection string
	db     *sql.DB
}

func (db *basedb) SetConnMaxLifetime(d time.Duration) {
	if db.db != nil {
		db.db.SetConnMaxLifetime(d)
	}
}

func (db *basedb) SetMaxIdleConns(n int) {
	if db.db != nil {
		db.db.SetMaxIdleConns(n)
	}
}

func (db *basedb) SetMaxOpenConns(n int) {
	if db.db != nil {
		db.db.SetMaxOpenConns(n)
	}
}

type Database interface {
	BindTable(tn string) Table
	CreateQuery(sql string) (Query, error)

	// no need to keep the two method?
	QueryRow(sql string, obj interface{}, args ...interface{}) error
	Query(sql string, objs interface{}, args ...interface{}) (err error)

	// if args[0] is a struct, use it's field
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Close() error

	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
}

type Table interface {
	Query
	Insert(obj interface{}) (sql.Result, error)
	// by pk
	Delete(obj interface{}) (sql.Result, error)
	Update(obj interface{}) (sql.Result, error)

	//Drop() error
	//Truncate() error
}

type Query interface {
	// need first call Exec
	Exec(args ...interface{}) error
	// after calling Exec, you can ethier Next nor All. Both Next and All will release the conn after the end.
	Next(obj interface{}, args ...interface{}) error
	All(objs interface{}) (err error)
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
	All(objs interface{}) (err error)

	//Count(filter interface{}) (int64, error)
	//Columns() ([]string, error)
	//Filter()
}
