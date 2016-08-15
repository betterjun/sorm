package sorm

import (
	"database/sql"
	"testing"
)

func printResult(t *testing.T, res sql.Result) {
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	t.Log("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
}

func TestFunction(t *testing.T) {
	db := NewDatabase("mysql", "root:root@tcp(127.0.0.1:3306)/world")
	if db == nil {
		t.Fatal("create db failed")
	}
	defer db.Close()

	sql := "DROP TABLE IF EXISTS xx"
	res, err := db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	printResult(t, res)

	sql = "CREATE TABLE xx(id int, name varchar(255), dummy varchar(255))"
	res, err = db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	printResult(t, res)

	sql = "INSERT INTO xx(id, name, dummy) VALUES(1, \"test_name\", \"dummy_string\")"
	res, err = db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	printResult(t, res)

	sql = "INSERT INTO xx(id, name, dummy) VALUES(2, \"test_name_2\", \"dummy_string_2\")"
	res, err = db.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	printResult(t, res)

	// select a single value
	sql = "select id from xx"
	var id int64
	err = db.QueryRow(sql, &id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("query a single value ok, id=%v\n", id)

	// select a map value
	var name string
	var dummy string
	mm := make(map[string]interface{})
	mm["id"] = &id
	mm["name"] = &name
	mm["dummy"] = &dummy
	sql = "select id, name from xx"
	err = db.QueryRow(sql, &mm)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("query a map value ok, id=%v, name=%v, dummy=%v\n", id, name, dummy)

	type tbs struct {
		SId   int `orm:"id"`
		Name  string
		Dummy string `orm:_`
	}
	r := &tbs{}

	// select a struct value
	sql = "select * from xx where id=?"
	err = db.QueryRow(sql, r, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*r)

	t.Log("db.Query")
	// query only supports a struct pointer for model register
	r = &tbs{}
	sql = "select * from xx where id>?"
	var allrows []tbs
	err = db.Query(sql, &allrows, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range allrows {
		t.Log(v)
	}
	//return

	si := make([]int, 0)
	sql = "select id from xx where id>?"
	err = db.Query(sql, &si, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range si {
		t.Log(v)
	}

	ss := make([]string, 0)
	sql = "select name from xx where id>?"
	err = db.Query(sql, &ss, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range ss {
		t.Log(v)
	}

	// test query
	sql = "select * from xx where id=?"
	q, err := db.CreateQuery(sql)
	if err != nil {
		t.Fatal(err)
	}

	err = q.ExecuteQuery(1)
	if err != nil {
		t.Fatal(err)
	}

	r = &tbs{}
	//t.Log("bb", *r)
	for q.Next(r) == nil {
		t.Log(*r)
	}

	err = q.ExecuteQuery(1)
	if err != nil {
		t.Fatal(err)
	}

	var all []tbs
	err = q.All(&all)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range all {
		t.Log(v)
	}
	all = nil

	// test table
	tb := db.BindTable("xx")
	res, err = tb.Insert(r)
	if err != nil {
		t.Fatal(err)
	}

	err = tb.ExecuteQuery()
	if err != nil {
		t.Fatal(err)
	}
	r = &tbs{}
	//t.Log("bb", *r)
	for tb.Next(r) == nil {
		t.Log(*r)
	}

	err = tb.ExecuteQuery()
	if err != nil {
		t.Fatal(err)
	}

	err = tb.All(&all)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range all {
		t.Log(v)
	}
	all = nil

	res, err = tb.Update(r)
	if err != nil {
		t.Fatal(err)
	}

	res, err = tb.Delete(r)
	if err != nil {
		t.Fatal(err)
	}

	res, err = tb.Insert(r)
	if err != nil {
		t.Fatal(err)
	}

}
