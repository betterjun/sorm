package sorm

import (
	"fmt"
	"testing"
	"time"
)

func dropTable(db Database, t *testing.T) {
	// test Exec without args
	sql := "DROP TABLE IF EXISTS xx"
	res, err := db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	if lii != 0 || ra != 0 {
		t.Fatalf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
	}
}

func createTable(db Database, t *testing.T) {
	sql := "CREATE TABLE xx(id int, name varchar(255), dummy varchar(255))"
	res, err := db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	if lii != 0 || ra != 0 {
		t.Fatalf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
	}
}

func insertTable(db Database, t *testing.T) {
	// test Exec with args
	count := 4
	for i := 1; i < count; i++ {
		sql := "INSERT INTO xx(id, name, dummy) VALUES(?,?,?)"
		res, err := db.Exec(sql, i, fmt.Sprintf("name%v", i), fmt.Sprintf("dummy%v", i))
		if err != nil {
			t.Fatal(err)
		}
		lii, _ := res.LastInsertId()
		ra, _ := res.RowsAffected()
		if lii != 0 || ra != 1 {
			t.Fatalf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
		}
	}
}

func TestDatabase(t *testing.T) {
	db := NewDatabase("mysql", CONN_STRING)
	if db == nil {
		t.Fatal("TestDatabase: create db failed")
	}

	db.SetConnMaxLifetime(time.Duration(10))
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	dropTable(db, t)
	createTable(db, t)
	insertTable(db, t)

	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}
}
