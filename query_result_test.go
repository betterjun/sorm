package sorm

import (
	"fmt"
	"io"
	"testing"
)

func testNext(db Database, t *testing.T) {
	// test CreateQuery
	sql := "select * from xx where id=?"
	q, err := db.CreateQuery(sql)
	if err != nil {
		t.Fatal(err)
	}

	// test Exec
	res, err := q.Exec(1)
	if err != nil {
		t.Fatal(err)
	}

	colNames, err := res.ColumnNames()
	if err != nil {
		t.Fatal(err)
	}
	if len(colNames) != 3 {
		t.Errorf("test Next ColumnNames failed, got %v records, expect 3\n", len(colNames))
	}
	names := []string{"id", "name", "dummy"}
	for i, name := range colNames {
		if name != names[i] {
			t.Errorf("test All ColumnNames failed, got %q, expect %q\n", name[0], names[i])
		}
	}

	// test Next struct
	count := 0
	r := &tbs{}
	for res.Next(r) == nil {
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
	res, err = q.Exec(2)
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	for {
		err = res.Next(&mm)
		if err != nil {
			if err == io.EOF {
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
	res, err = q.Exec(3)
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	id = 0
	name = ""
	dummy = ""
	for {
		err = res.Next(&id, &name, &dummy)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		if id != 3 {
			t.Errorf("test Next built-in type failed, id=%v, expect 3\n", id)
		}
		if name != "name3" {
			t.Errorf("test Next built-in type failed, name=%v, expect \"name3\"\n", name)
		}
		if dummy != "dummy3" {
			t.Errorf("test Next built-in type failed, dummy=%v, expect \"dummy3\"\n", dummy)
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
}

func testAll(db Database, t *testing.T) {
	sql := "select id from xx where id>? and id<? order by id asc"
	q, err := db.CreateQuery(sql)
	if err != nil {
		t.Fatal(err)
	}

	// test All built-in type
	res, err := q.Exec(0, 10)
	if err != nil {
		t.Fatal(err)
	}

	colNames, err := res.ColumnNames()
	if err != nil {
		t.Fatal(err)
	}
	if len(colNames) != 1 {
		t.Errorf("test All ColumnNames failed, got %v records, expect 1\n", len(colNames))
	}
	if colNames[0] != "id" {
		t.Errorf("test All ColumnNames failed, got %q, expect %q\n", colNames[0], "id")
	}

	var si []int
	err = res.All(&si)
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

	res, err = q.Exec(0, 10)
	if err != nil {
		t.Fatal(err)
	}

	var allrows []tbs
	err = res.All(&allrows)
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

func TestQuery(t *testing.T) {
	db := NewDatabase("mysql", CONN_STRING)
	if db == nil {
		t.Fatal("TestDatabase: create db failed")
	}
	defer db.Close()

	testNext(db, t)
	testAll(db, t)
}
