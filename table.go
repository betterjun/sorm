package sorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type table struct {
	name string
	db   Database
}

func (t *table) Insert(values ...interface{}) (res sql.Result, err error) {
	if t.db == nil {
		return nil, fmt.Errorf("db is not opened")
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("Table.Insert must have an input value")
	}

	insertSql := fmt.Sprintf("insert into %v(", t.name)
	valueSql := "("
	args := make([]interface{}, 0)
	var sql string

	obj := values[0]
	obv := reflect.ValueOf(obj)
	if obv.Kind() == reflect.Ptr {
		obv = obv.Elem()
	}

	switch obv.Kind() {
	case reflect.Struct:
		tis := getFieldInfoFromStruct(obv)
		if len(tis) == 0 {
			return
		}
		for _, v := range tis {
			if v.fn == "_" {
				continue
			}

			insertSql += v.fn + ","
			valueSql += fmt.Sprintf("?,")
			args = append(args, v.fp.Interface())
		}

		sql = insertSql[0:len(insertSql)-1] + ") values" + valueSql[0:len(valueSql)-1] + ")"
	case reflect.Map:
		keys := obv.MapKeys()
		for _, v := range keys {
			k := v.Interface().(string)
			insertSql += k + ","
			valueSql += fmt.Sprintf("?,")
			args = append(args, obv.MapIndex(v).Interface())
		}

		sql = insertSql[0:len(insertSql)-1] + ") values" + valueSql[0:len(valueSql)-1] + ")"
	default:
		insertSql = fmt.Sprintf("insert into %v values(", t.name)
		insertSql += strings.Repeat("?,", len(values))
		args = values
		sql = insertSql[0:len(insertSql)-1] + ")"
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no valid fields found in the object")
	}
	fmt.Println("table.Insert", sql)
	return t.db.Exec(sql, args...)
}

func (t *table) Delete(filter string) (res sql.Result, err error) {
	if t.db == nil {
		return nil, fmt.Errorf("db is not opened")
	}

	var sql string
	if filter == "" {
		sql = "delete from " + t.name
	} else {
		sql = "delete from " + t.name + " where " + filter
	}

	return t.db.Exec(sql)
}

func (t *table) Update(filter string, value interface{}) (res sql.Result, err error) {
	if t.db == nil {
		return nil, fmt.Errorf("db is not opened")
	}

	updateSql := fmt.Sprintf("update %v set ", t.name)
	whereSql := ""
	args := make([]interface{}, 0)

	obv := reflect.ValueOf(value)
	if obv.Kind() == reflect.Ptr {
		obv = obv.Elem()
	}

	switch obv.Kind() {
	case reflect.Struct:
		tis := getFieldInfoFromStruct(obv)
		if len(tis) == 0 {
			return nil, fmt.Errorf("no receiver fields found")
		}
		for _, v := range tis {
			if v.fn == "_" {
				continue
			}

			whereSql += v.fn + "=?,"
			args = append(args, v.fp.Interface())
		}
	case reflect.Map:
		keys := obv.MapKeys()
		for _, v := range keys {
			k := v.Interface().(string)
			whereSql += k + "=?,"
			args = append(args, obv.MapIndex(v).Interface())
		}
	default:
		return nil, fmt.Errorf("non-supported input args, only supporting struct or map[string]interface{}")
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no valid fields found in the object")
	}

	var sql string
	if filter == "" {
		sql = updateSql + whereSql[0:len(whereSql)-1]
	} else {
		sql = updateSql + whereSql[0:len(whereSql)-1] + " where " + filter
	}
	fmt.Println("table.Update", sql)
	return t.db.Exec(sql, args...)
}

func (t *table) Query(filter string) (res Result, err error) {
	if t.db == nil {
		return nil, fmt.Errorf("db is not opened")
	}

	var sql string
	if filter == "" {
		sql = "select * from " + t.name
	} else {
		sql = "select * from " + t.name + " where " + filter
	}
	q, err := t.db.CreateQuery(sql)
	if err != nil {
		return nil, err
	}

	return q.Exec()
}
