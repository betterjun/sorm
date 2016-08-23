package sorm

import (
	"fmt"
	"io"
	"testing"
)

type helper func(res Result, t *testing.T)

func testQueryHelper(tb Table, t *testing.T, hf helper) {
	filter := ""
	res, err := tb.Query(filter)
	if err != nil {
		t.Fatal(err)
	}

	hf(res, t)
}

func testTableQuery(tb Table, t *testing.T) {
	testQueryHelper(tb, t, testNextBuiltin)
	testQueryHelper(tb, t, testNextMap)
	testQueryHelper(tb, t, testNextStruct)

	//testQueryHelper(tb, t, testAllBuiltin)
	testQueryHelper(tb, t, testAllStruct)
}

func testNextBuiltin(tb Result, t *testing.T) {
	// test Next built-in type
	var err error
	count := 0
	var id int
	var name string
	var dummy string
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

func testNextMap(tb Result, t *testing.T) {
	// test Next map[string]interface{}
	var err error
	count := 0
	var id int
	var name string
	var dummy string
	mm := make(map[string]interface{})
	mm["id"] = &id
	mm["name"] = &name
	mm["dummy"] = &dummy

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
}

func testNextStruct(tb Result, t *testing.T) {
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

	err := tb.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func testAllBuiltin(tb Result, t *testing.T) {
	// test All built-in type
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
}

func testAllStruct(tb Result, t *testing.T) {
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
	// test insert built-in type, you should passing all the columns by it's order in db.
	id := 1000
	name := "name1001"
	dummy := "dummy1001"
	res, err := tb.Insert(id, name, dummy)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	// test insert map
	id++
	args := make(map[string]interface{})
	args["id"] = id
	args["name"] = "insert by map"
	res, err = tb.Insert(&args)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	// test insert struct
	id++
	ts := &tbs{SId: id, Name: "fn=name", Dummy: "ignored"}
	res, err = tb.Insert(ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
}

func testTableUpdate(tb Table, t *testing.T) {
	// test update map
	id := 1000
	args := make(map[string]interface{})
	args["name"] = "updated by map"
	args["dummy"] = "updated by testTableUpdate"

	filter := "id>=1000"
	res, err := tb.Update(filter, args)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)

	// test update struct
	id++
	ts := &tbs{SId: id}
	ts.Name = "ts"
	ts.Dummy = "dummy"
	filter = fmt.Sprintf("id=%v and name <>\"\"", id)
	res, err = tb.Update(filter, ts)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ = res.LastInsertId()
	ra, _ = res.RowsAffected()
	t.Logf("res.LastInsertId()=%v, res.RowsAffected()=%v", lii, ra)
}

func testTableDelete(tb Table, t *testing.T) {
	// test delete
	filter := "id>=1000"
	res, err := tb.Delete(filter)
	if err != nil {
		t.Fatal(err)
	}
	lii, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
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

	testTableQuery(tb, t)
	testTableInsert(tb, t)
	testTableUpdate(tb, t)
	testTableDelete(tb, t)
}
