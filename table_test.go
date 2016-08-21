package sorm

import (
	"fmt"
	"io"
	"testing"
)

func testTableNext(tb Table, t *testing.T) {
	// test Next struct
	count := 0
	r := &tbs{}
	for tb.Next(r) == nil {
		count++
		if r.SId != count {
			t.Errorf("test Next struct failed, r.SId=%v, expect %v\n", r.SId, count)
		}
		if r.Name != "" {
			t.Errorf("test Next struct failed, name=%q, expect %q\n", r.Name, "")
		}
		if r.Dummy != fmt.Sprintf("dummy%v", count) {
			t.Errorf("test Next struct failed, dummy=%q, expect %q\n", r.Dummy, fmt.Sprintf("dummy%v", count))
		}
	}
	if count != 3 {
		t.Errorf("test Next struct failed, got %v records, expect 3\n", count)
	}
	//var err error
	err := tb.Close()
	if err != nil {
		t.Fatal(err)
	}

	// test Next map[string]interface{}
	var id int
	var name string
	var dummy string
	mm := make(map[string]interface{})
	mm["id"] = &id
	mm["name"] = &name
	mm["dummy"] = &dummy

	count = 0
	for {
		err = tb.Next(&mm)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		count++

		if id != count {
			t.Errorf("test Next map[string]interface{} failed, id=%v, expect %v\n", id, count)
		}
		if name != fmt.Sprintf("name%v", count) {
			t.Errorf("test Next map[string]interface{} failed, name=%q, expect %q\n", name, fmt.Sprintf("name%v", count))
		}
		if dummy != fmt.Sprintf("dummy%v", count) {
			t.Errorf("test Next map[string]interface{} failed, dummy=%q, expect %q\n", dummy, fmt.Sprintf("dummy%v", count))
		}
	}
	if count != 3 {
		t.Errorf("test Next map[string]interface{} failed, got %v records, expect 3\n", count)
	}
	err = tb.Close()
	if err != nil {
		t.Fatal(err)
	}

	// test Next built-in type
	count = 0
	id = 0
	name = ""
	dummy = ""
	for {
		err = tb.Next(&id, &name, &dummy)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		count++
		if id != count {
			t.Errorf("test Next built-in type failed, id=%v, expect %v\n", id, count)
		}
		if name != fmt.Sprintf("name%v", count) {
			t.Errorf("test Next built-in type failed, name=%v, expect %q\n", name, fmt.Sprintf("name%v", count))
		}
		if dummy != fmt.Sprintf("dummy%v", count) {
			t.Errorf("test Next built-in type failed, dummy=%v, expect %q\n", dummy, fmt.Sprintf("dummy%v", count))
		}
	}
	if count != 3 {
		t.Errorf("test Next built-in type failed, got %v records, expect 3\n", count)
	}

	err = tb.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func testTableAll(tb Table, t *testing.T) {
	// test All built-in type

	/*
		var si []int
		err := tb.All(&si)
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
		err = tb.Close()
		if err != nil {
			t.Fatal(err)
		}
	*/
	// test All struct
	var allrows []tbs
	err := tb.All(&allrows)
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
	err = tb.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func testTableInsert(tb Table, t *testing.T) {
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
}

func testTableUpdate(tb Table, t *testing.T) {
	ts := &tbs{SId: 1000}
	ts.Name = "ts"
	res, err := tb.Update(ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
}

func testTableDelete(tb Table, t *testing.T) {
	ts := &tbs{SId: 1000}
	res, err := tb.Delete(ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	ts.SId = 2000
	res, err = tb.Delete(ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
}

func TestTable(t *testing.T) {
	db := NewDatabase("mysql", CONN_STRING)
	if db == nil {
		t.Fatal("TestDatabase: create db failed")
	}
	defer db.Close()

	// test table
	tb, err := db.BindTable("xx")

	if err != nil {
		t.Fatal(err)
	}

	testTableNext(tb, t)
	testTableAll(tb, t)
	testTableInsert(tb, t)
	testTableUpdate(tb, t)
	testTableDelete(tb, t)
}
