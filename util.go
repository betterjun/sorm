package sorm

import (
	"database/sql"
	"reflect"
	"strconv"
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
		ind := reflect.Indirect(v) // equal with v.Elem()
		switch ind.Kind() {
		case reflect.Map:
			return getScanFieldFromMap(ind, cols)
		case reflect.Struct:
			return getScanFieldFromStruct(ind, cols)
		default: // pointer to value
			scanArgs = append(scanArgs, ptr)
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
	return getFields(getFieldsFromStruct(v), cols)
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

func getFieldsFromStruct(v reflect.Value) (fields map[string]interface{}) {
	fields = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("orm")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		//fmt.Println("name =", name, ",val =", v.Field(i).Addr(), v.Field(i).Addr().Elem(), v.Field(i).Addr().Elem().CanAddr())
		// take the addr and as interface{}, this is required by sql Scan
		fields[name] = v.Field(i).Addr().Interface()
	}

	return fields
}

type tagInfo struct {
	fn      string
	pk      bool
	ignored bool
	fp      reflect.Value
}

func getFieldInfoFromStruct(v reflect.Value) (fields map[string]*tagInfo) {
	fields = make(map[string]*tagInfo)
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag

		ti := parseTag(strings.ToLower(fieldInfo.Name), tag.Get("orm"))
		ti.fp = v.Field(i).Addr()
		if ti != nil && ti.fn != "_" {
			fields[ti.fn] = ti
		}
	}

	return fields
}

func parseTag(fieldName, tag string) (ti *tagInfo) {
	/*
		`orm:"pk=true"`
		`orm:"_"`
		`orm:"pk=false;fn=name"`
	*/
	fieldName = strings.ToLower(fieldName)

	ti = &tagInfo{fn: fieldName, pk: false, ignored: false}

	tags := strings.Split(tag, ";")
	if len(tags) > 0 {
		for _, kvp := range tags {
			kvp = strings.TrimSpace(kvp)
			if kvp == "_" {
				ti.fn = "_"
				ti.pk = false
				ti.ignored = true
				continue
			}

			kv := strings.Split(kvp, "=")
			if len(kv) != 2 { // wrong format of orm, just use fieldname
				continue
			}
			kv[0] = strings.TrimSpace(kv[0])
			kv[1] = strings.TrimSpace(kv[1])

			if kv[0] == "pk" {
				ti.pk, _ = strconv.ParseBool(kv[1])
			} else if kv[0] == "fn" {
				if kv[1] == "_" {
					ti.fn = "_"
					ti.pk = false
					ti.ignored = true
				} else if kv[1] == "" {
					ti.fn = fieldName
				} else {
					ti.fn = kv[1]
				}
			}
		}
	}

	return
}
