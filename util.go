package sorm

import (
	"database/sql"
	"reflect"
	"strings"
)

func getFieldsForOne(ptr interface{}, optPtr []interface{}, cols []string) (scanArgs []interface{}) {
	v := reflect.ValueOf(ptr)
	switch v.Kind() {
	case reflect.Ptr: // only accept pointer
		ind := reflect.Indirect(v) // equal with v.Elem()
		switch ind.Kind() {
		case reflect.Map:
			return getScanFieldFromMap(ind, cols)
		case reflect.Struct:
			return getScanFieldFromStruct(ind, cols)
		default: // pointer to value
			scanArgs = append(scanArgs, ptr)
			for _, op := range optPtr {
				if reflect.ValueOf(op).Kind() != reflect.Ptr {
					return nil
				}
				scanArgs = append(scanArgs, op)
			}
			return scanArgs
		}
	default:
		return nil
	}
}

func getScanFieldFromMap(v reflect.Value, cols []string) (scanArgs []interface{}) {
	fields := make(map[string]interface{})
	for _, k := range v.MapKeys() {
		fields[k.Interface().(string)] = v.MapIndex(k).Interface()
	}
	return getFields(fields, cols)
}

func getScanFieldFromStruct(v reflect.Value, cols []string) (scanArgs []interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		ti := parseTag(fieldInfo.Name, fieldInfo.Tag.Get("sorm"))
		if ti != nil && ti.fn != "_" {
			fields[ti.fn] = v.Field(i).Addr().Interface()
		}
	}
	return getFields(fields, cols)
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

// name, value pointer for a struct field.
type tagInfo struct {
	fn string
	fp reflect.Value
}

func getFieldInfoFromStruct(v reflect.Value) (fields map[string]*tagInfo) {
	fields = make(map[string]*tagInfo)
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag

		ti := parseTag(fieldInfo.Name, tag.Get("sorm"))
		ti.fp = v.Field(i).Addr()
		if ti != nil && ti.fn != "_" {
			fields[ti.fn] = ti
		}
	}
	return fields
}

/*
supported tag:
	`sorm:"_"`
	`sorm:"fn=name"`
*/
func parseTag(fieldName, tag string) (ti *tagInfo) {
	fieldName = strings.ToLower(fieldName)
	ti = &tagInfo{fn: fieldName}
	tags := strings.Split(tag, ";")
	if len(tags) > 0 {
		for _, kvp := range tags {
			kvp = strings.TrimSpace(kvp)
			if kvp == "_" {
				ti.fn = "_"
				continue
			}

			kv := strings.Split(kvp, "=")
			if len(kv) != 2 { // wrong format of orm, just use fieldname
				continue
			}
			kv[0] = strings.TrimSpace(kv[0])
			kv[1] = strings.TrimSpace(kv[1])

			if kv[0] == "fn" {
				if kv[1] == "_" {
					ti.fn = "_"
				} else if kv[1] == "" {
					ti.fn = fieldName
				} else {
					ti.fn = kv[1]
				}
			}
		}
	}

	return ti
}
