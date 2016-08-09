package sorm

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//
func NewDatabase(dbtype, conn string) (db Database) {
	d := &database{}
	err := d.open(conn)
	if err != nil {
		fmt.Printf("sorm create database connection failed:%v\n", err)
		return nil
	}
	return d
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

func (t *table) Query(objs []interface{}) (err error) {
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

func (q *query) Query(objs []interface{}) (err error) {
	return
}

type database struct {
	dbtype string // just now only support "mysql"
	dsn    string // connection string
	db     *sql.DB
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

func (db *database) QueryRow(sql string, obj interface{}, args ...interface{}) (err error) {
	return err
}
func (db *database) Query(sql string, objs []interface{}, args ...interface{}) (err error) {
	return err
}

func (db *database) Exec(sql string, args ...interface{}) (res sql.Result, err error) {
	if db.db != nil {
		err = db.db.Close()
		if err == nil {
			db.db = nil
		}
	}
	return res, err
}
