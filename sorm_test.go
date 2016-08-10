package sorm

import "testing"

func TestFunction(t *testing.T) {
	db := NewDatabase("mysql", "root:root@tcp(127.0.0.1:3306)/world")
	if db == nil {
		t.Fatal("create db failed")
	}
	defer db.Close()

	res, err := db.Exec("create table xx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)

	sql := "select id from xx"
	var id int64
	err = db.QueryRow(sql, &id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	var name string
	mm := make(map[string]interface{})
	mm["id"] = &id
	mm["name"] = &name
	sql = "select id, name from xx"
	err = db.QueryRow(sql, &mm)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
	t.Log(name)

	type tbs struct {
		SId  int `orm:"id"`
		Name string
	}
	r := &tbs{}
	err = db.QueryRow(sql, r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*r)

	all, err := db.Query(sql, r)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range all {
		t.Log(v)
	}

	tb := db.BindTable("xx")
	res, err = tb.Insert(r)
	if err != nil {
		t.Fatal(err)
	}

	res, err = tb.Update(r)
	if err != nil {
		t.Fatal(err)
	}

	res, err = tb.Delete(r)
	if err != nil {
		t.Fatal(err)
	}

	err = tb.QueryRow(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*r)

	res, err = tb.Insert(r)
	if err != nil {
		t.Fatal(err)
	}

	all, err = tb.Query(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range all {
		t.Log(v)
	}

	q, err := db.CreateQuery("select * from xx")
	if err != nil {
		t.Fatal(err)
	}

	res, err = q.Exec()
	if err != nil {
		t.Fatal(err)
	}

	err = q.QueryRow(r)
	if err != nil {
		t.Fatal(err)
	}

	all, err = q.Query(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range all {
		t.Log(v)
	}

	for q.Next(r) != nil {
		t.Log(*r)
	}

}
