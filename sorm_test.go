package sorm

import (
	"fmt"
	"testing"
	"time"
)

type tbs struct {
	SId   int    `orm:"pk=1;fn=id"`
	Name  string `orm:"_"`
	Dummy string `orm:"fn=dummy"`
}

const (
	CONN_STRING = "root:root@tcp(127.0.0.1:3306)/world"
)

func TestDatabase(t *testing.T) {
	db := NewDatabase("mysql", CONN_STRING)
	if db == nil {
		t.Fatal("TestDatabase: create db failed")
	}

	// test Exec without args
	sql := "DROP TABLE IF EXISTS xx"
	res, err := db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	sql = "CREATE TABLE xx(id int, name varchar(255), dummy varchar(255))"
	res, err = db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	// test Exec with args
	count := 4
	for i := 1; i < count; i++ {
		sql = "INSERT INTO xx(id, name, dummy) VALUES(?,?,?)"
		res, err = db.Exec(sql, i, fmt.Sprintf("name%v", i), fmt.Sprintf("dummy%v", i))
		if err != nil {
			t.Fatal(err)
		}
		lii, _ = res.LastInsertId()
		ra, _ = res.RowsAffected()
		t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
	}

	// test QueryRow, built-in type pointer
	sql = "select name from xx where id=?"
	var name string
	err = db.QueryRow(sql, &name, 1)
	if err != nil {
		t.Fatal(err)
	}
	if name != "name1" {
		t.Errorf("test QueryRow, built-in type pointer failed, name=%v, expect \"name1\"\n", name)
	}
	t.Logf("test QueryRow, built-in type pointer ok")

	// test QueryRow, map of pointers
	name = ""
	var id int64
	var dummy string
	mm := make(map[string]interface{})
	mm["id"] = &id
	mm["name"] = &name
	mm["dummy"] = &dummy
	sql = "select id, name from xx where id = ?"
	err = db.QueryRow(sql, &mm, 2)
	if err != nil {
		t.Fatal(err)
	}
	if id != 2 {
		t.Errorf("test QueryRow, map of pointers failed, id=%v, expect 2\n", id)
	}
	if name != "name2" {
		t.Errorf("test QueryRow, map of pointers failed, name=%q, expect \"name2\"\n", name)
	}
	if dummy != "" {
		t.Errorf("test QueryRow, map of pointers failed, dummy=%q, expect \"\"\n", dummy)
	}
	t.Logf("test QueryRow, map of pointers ok")

	r := &tbs{}

	// test QueryRow, a struct pointer
	sql = "select * from xx where id=?"
	err = db.QueryRow(sql, r, 3)
	if err != nil {
		t.Fatal(err)
	}
	if r.SId != 3 {
		t.Errorf("test QueryRow, a struct pointer failed, r.SId=%v, expect 3\n", r.SId)
	}
	if r.Name != "" {
		t.Errorf("test QueryRow, a struct pointer failed, r.Name=%q, expect \"\"\n", r.Name)
	}
	if r.Dummy != "dummy3" {
		t.Errorf("test QueryRow, a struct pointer failed, r.Dummy=%q, expect \"dummy3\"\n", r.Dummy)
	}
	t.Logf("test QueryRow, a struct pointer ok")

	// test Query, a struct slice
	r = &tbs{}
	sql = "select * from xx where id>? and id<? order by id asc"
	var allrows []tbs
	err = db.Query(sql, &allrows, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(allrows) != 3 {
		t.Errorf("test Query, a struct slice failed, slice length=%v, expect 3\n", len(allrows))
	}
	for i, r := range allrows {
		if r.SId != i+1 {
			t.Errorf("test Query, a struct slice failed, r.SId=%v, expect %v\n", r.SId, i+1)
		}
		if r.Name != "" {
			t.Errorf("test Query, a struct slice failed, r.Name=%q, expect \"\"\n", r.Name)
		}
		if r.Dummy != fmt.Sprintf("dummy%v", i+1) {
			t.Errorf("test Query, a struct slice failed, r.Dummy=%q, expect %q\n", r.Dummy, fmt.Sprintf("dummy%v", i+1))
		}
	}
	t.Logf("test Query, a struct slice ok")

	// test Query, a built-in type slice
	si := make([]int, 0)
	sql = "select id from xx where id>?"
	err = db.Query(sql, &si, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range si {
		if v != i+1 {
			t.Errorf("test Query, a built-in type slice failed, id=%v, expect %v\n", v, i+1)
		}
	}
	t.Logf("test Query, a built-in type slice")

	db.SetConnMaxLifetime(time.Duration(10))
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	db := NewDatabase("mysql", CONN_STRING)
	if db == nil {
		t.Fatal("TestDatabase: create db failed")
	}
	defer db.Close()

	// test CreateQuery
	sql := "select * from xx where id=?"
	q, err := db.CreateQuery(sql)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("CreateQuery ok")

	// test ExecuteQuery
	err = q.ExecuteQuery(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ExecuteQuery ok")

	// test Next struct
	count := 0
	r := &tbs{}
	for q.Next(r) == nil {
		if r.SId != 1 {
			t.Errorf("test Next struct failed, r.SId=%v, expect 1\n", r.SId)
		}
		count++
	}
	if count != 1 {
		t.Errorf("test Next struct failed, got %v records, expect 1\n", count)
	}

	// test Next map[string]interface{}
	var id int64
	var name string
	var dummy string
	mm := make(map[string]interface{})
	mm["id"] = &id
	mm["name"] = &name
	mm["dummy"] = &dummy
	err = q.ExecuteQuery(2)
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	for {
		err = q.Next(&mm)
		if err != nil {
			if err.Error() == "end of query results" {
				break
			}
			t.Fatal(err)
		}
		if id != 2 {
			t.Errorf("test Next map[string]interface{} failed, id=%v, expect 2\n", id)
		}
		if name != "name2" {
			t.Errorf("test Next map[string]interface{} failed, name=%v, expect \"name2\"\n", name)
		}
		if dummy != "dummy2" {
			t.Errorf("test Next map[string]interface{} failed, dummy=%v, expect \"dummy2\"\n", dummy)
		}
		count++
	}
	if count != 1 {
		t.Errorf("test Next map[string]interface{} failed, got %v records, expect 1\n", count)
	}

	// test Next built-in type
	err = q.ExecuteQuery(3)
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	id = 0
	for {
		err = q.Next(&id)
		if err != nil {
			if err.Error() == "end of query results" {
				break
			}
			t.Fatal(err)
		}
		if id != 3 {
			t.Errorf("test Next built-in type failed, id=%v, expect 3\n", id)
		}
		count++
	}
	if count != 1 {
		t.Errorf("test Next built-in type failed, got %v records, expect 1\n", count)
	}

	err = q.Close()
	if err != nil {
		t.Fatal(err)
	}

	sql = "select id from xx where id>? and id<? order by id asc"
	q, err = db.CreateQuery(sql)
	if err != nil {
		t.Fatal(err)
	}

	// test All built-in type
	err = q.ExecuteQuery(0, 10)
	if err != nil {
		t.Fatal(err)
	}

	var si []int
	err = q.All(&si)
	if err != nil {
		t.Fatal(err)
	}
	if len(si) != 3 {
		t.Errorf("test All built-in type failed, slice length=%v, expect 3\n", len(si))
	}
	for i, r := range si {
		if r != i+1 {
			t.Errorf("test All built-in type failed, r=%v, expect %v\n", r, i+1)
		}
	}

	// test All struct
	sql = "select * from xx where id>? and id<? order by id asc"
	q, err = db.CreateQuery(sql)
	if err != nil {
		t.Fatal(err)
	}

	err = q.ExecuteQuery(0, 10)
	if err != nil {
		t.Fatal(err)
	}

	var allrows []tbs
	err = q.All(&allrows)
	if err != nil {
		t.Fatal(err)
	}
	if len(allrows) != 3 {
		t.Errorf("test All struct failed, slice length=%v, expect 3\n", len(allrows))
	}
	for i, r := range allrows {
		if r.SId != i+1 {
			t.Errorf("test All struct failed, r.SId=%v, expect %v\n", r.SId, i+1)
		}
		if r.Name != "" {
			t.Errorf("test All struct failed, r.Name=%q, expect \"\"\n", r.Name)
		}
		if r.Dummy != fmt.Sprintf("dummy%v", i+1) {
			t.Errorf("test All struct failed, r.Dummy=%q, expect %q\n", r.Dummy, fmt.Sprintf("dummy%v", i+1))
		}
	}
}

func TestTable(t *testing.T) {
	db := NewDatabase("mysql", CONN_STRING)
	if db == nil {
		t.Fatal("TestDatabase: create db failed")
	}
	defer db.Close()

	// test table
	tb := db.BindTable("xx")

	err := tb.ExecuteQuery()
	if err != nil {
		t.Fatal(err)
	}
	r := &tbs{}

	for tb.Next(r) == nil {
		t.Log(*r)
	}

	err = tb.ExecuteQuery()
	if err != nil {
		t.Fatal(err)
	}

	allrows := make([]tbs, 0)
	err = tb.All(&allrows)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range allrows {
		t.Log(v)
	}
	allrows = nil

	ts := &tbs{SId: 1000, Name: "fn=name", Dummy: "ignored"}
	res, err := tb.Insert(ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	args := make(map[string]interface{})
	args["id"] = 2000
	args["name"] = "inserte by map"
	res, err = tb.Insert(&args)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	ts.Name = "ts"
	res, err = tb.Update(ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	res, err = tb.Delete(ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	res, err = tb.Delete(&args)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

}
