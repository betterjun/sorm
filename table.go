package sorm

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type table struct {
	name string
	db   Database
}

func (t *table) Insert(values ...interface{}) (res sql.Result, err error) {
	// obj 可以为对象
	// 也可以为map[string]interface{}
	if t.db != nil {
		insertSql := fmt.Sprintf("insert into %v(", t.name)
		valueSql := "("

		obj := values[0]
		ro := reflect.ValueOf(obj)
		if ro.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("only support pointer argument")
		}

		obv := ro.Elem()

		if obv.Kind() == reflect.Struct {
			tis := getFieldInfoFromStruct(obv)
			if len(tis) == 0 {
				return
			}
			args := make([]interface{}, 0)
			for _, v := range tis {
				if v.fn == "_" {
					continue
				}

				insertSql += v.fn + ","
				valueSql += fmt.Sprintf("?,")
				args = append(args, v.fp.Interface())
			}

			if len(args) == 0 {
				return nil, fmt.Errorf("no valid fields found in the object")
			}

			sql := insertSql[0:len(insertSql)-1] + ") values" + valueSql[0:len(valueSql)-1] + ")"
			fmt.Println("table.Insert", sql)
			return t.db.Exec(sql, args...)
		} else if obv.Kind() == reflect.Map {
			keys := obv.MapKeys()
			args := make([]interface{}, 0)
			for _, v := range keys {
				k := v.Interface().(string)

				insertSql += k + ","
				valueSql += fmt.Sprintf("?,")
				args = append(args, obv.MapIndex(v).Interface())
			}

			if len(args) == 0 {
				return nil, fmt.Errorf("no valid fields found in the object")
			}

			sql := insertSql[0:len(insertSql)-1] + ") values" + valueSql[0:len(valueSql)-1] + ")"
			fmt.Println("table.Insert", sql)
			return t.db.Exec(sql, args...)
		} else {
			return nil, fmt.Errorf("non-supported input args, only supporting struct or map[string]interface{}")
		}
	}

	return nil, fmt.Errorf("db is not opened")
}

func (t *table) Delete(filter string) (res sql.Result, err error) {
	if t.db != nil {
		var sql string
		if filter == "" {
			sql = "delete from " + t.name
		} else {
			sql = "delete from " + t.name + " where " + filter
		}

		res, err = t.db.Exec(sql)
	}

	return nil, fmt.Errorf("db is not opened")
}

func (t *table) Update(filter string, value interface{}) (res sql.Result, err error) {
	// obj 可以为对象，根据pk来更新
	// Update(key, value) key,value都是map[string]interface{}

	if t.db != nil {
		ro := reflect.ValueOf(value)
		if ro.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("only support pointer argument")
		}

		obv := ro.Elem()
		tis := getFieldInfoFromStruct(obv)
		updateSql := fmt.Sprintf("update %v set ", t.name)
		whereSql := ""
		if len(tis) == 0 {
			return
		}

		args := make([]interface{}, 0)
		var whereArgs interface{}
		for _, v := range tis {
			if v.fn == "_" {
				continue
			} else if v.pk {
				whereSql += fmt.Sprintf(" where %v=?", v.fn)
				whereArgs = v.fp.Interface()
			} else {
				updateSql += v.fn + "=?, "
				args = append(args, v.fp.Interface())
			}
		}

		if whereArgs == nil {
			return nil, fmt.Errorf("no pk field found in the object")
		}
		if len(args) == 0 {
			return nil, fmt.Errorf("no valid fields found in the object")
		}
		args = append(args, whereArgs)

		sql := updateSql[0:len(updateSql)-2] + whereSql
		fmt.Println("table.Update", sql)
		return t.db.Exec(sql, args...)
	}

	return nil, fmt.Errorf("db is not opened")
}

func (t *table) Query(filter string) (res Result, err error) {
	if t.db != nil {
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

	return nil, fmt.Errorf("db is not opened")
}
