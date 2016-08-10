package sorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func getScanFields(ptr interface{}, cols []string) (scanArgs []interface{}) {
	/*
		如果ptr是指针，则取第一个
		如果ptr是map[string]interface{}，则取对应的值做参数,interface{}必须也为指针
		如果ptr是struct，则取对应的字段做参数,
		其他视为无效
	*/
	v := reflect.ValueOf(ptr)
	switch v.Kind() {
	case reflect.Ptr: // only accept pointer
		fmt.Println("s", 1)
		ind := reflect.Indirect(v) // equal with v.Elem()
		switch ind.Kind() {
		case reflect.Map:
			fmt.Println("s", 2)
			return getScanFieldFromMap(ind, cols)
		case reflect.Struct:
			fmt.Println("s", 4)
			return getScanFieldFromStruct(ind, cols)
		default: // pointer to value
			fmt.Println("s", 5, ind)
			scanArgs = append(scanArgs, ptr)
			return scanArgs
		}
	case reflect.Slice: // every element in slice should be pointer
		fmt.Println("s", 3)
		return getScanFieldFromSlice(v, cols)
	case reflect.Struct:
		fmt.Println("s", 6)
		return getScanFieldFromStruct(v, cols)
	default:
		return nil
	}
}

func getScanFieldFromMap(v reflect.Value, cols []string) (scanArgs []interface{}) {
	fields := make(map[string]interface{})
	for _, k := range v.MapKeys() {
		//fmt.Println("mp =", v.MapIndex(k))
		fields[k.Interface().(string)] = v.MapIndex(k).Interface()
	}
	return getFields(fields, cols)
}

func getScanFieldFromStruct(v reflect.Value, cols []string) (scanArgs []interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("orm")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fmt.Println("name =", name)
		// take the addr and as interface{}, this is required by sql Scan
		fields[name] = v.Field(i).Addr().Interface()
	}
	return getFields(fields, cols)
}

func getScanFieldFromSlice(v reflect.Value, cols []string) (scanArgs []interface{}) {
	vals := v.Interface().([]interface{})

	lc := len(cols)
	c := 0
	for _, vp := range vals {
		obj := reflect.ValueOf(vp).Elem() // the struct variable
		if obj.Kind() != reflect.Ptr {    // only accept pointer
			return nil
		}

		scanArgs = append(scanArgs, vp)
		c++
		if c > lc {
			break
		}
	}

	for ; c < lc; c++ {
		scanArgs = append(scanArgs, new(sql.RawBytes))
	}
	return scanArgs
}

func getFields(fields map[string]interface{}, cols []string) (scanArgs []interface{}) {
	for _, name := range cols {
		f := fields[name]
		if f == nil { // no receiver found in the struct, use a raw bytes to receive
			f = new(sql.RawBytes)
		}
		scanArgs = append(scanArgs, f)
	}
	return scanArgs
}
