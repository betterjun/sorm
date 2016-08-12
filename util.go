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
		fmt.Println("name =", name, ",val =", v.Field(i).Addr(), v.Field(i).Addr().Elem(), v.Field(i).Addr().Elem().CanAddr())
		// take the addr and as interface{}, this is required by sql Scan
		//fields[name] = v.Field(i).Addr().Interface()
		if v.Field(i).CanAddr() {
			fmt.Printf("%v CanAddr\n", name)
			fields[name] = v.Field(i).Addr().Interface()
		} else {
			fmt.Printf("%v Can not Addr\n", name)
			fields[name] = v.Field(i).Addr().Interface()
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

func setFieldValue(ind reflect.Value, value interface{}) {
	switch ind.Kind() {
	case reflect.Bool:
		if v, ok := value.(bool); ok {
			ind.SetBool(v)
		} else {
			ind.SetBool(false)
		}

	case reflect.String:
		if value == nil {
			ind.SetString("")
		} else {
			ind.SetString(toString(value))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			ind.SetInt(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ind.SetInt(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ind.SetInt(int64(val.Uint()))
			default:
				ind.SetInt(0)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			ind.SetUint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ind.SetUint(uint64(val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ind.SetUint(val.Uint())
			default:
				ind.SetUint(0)
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			ind.SetFloat(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				ind.SetFloat(val.Float())
			default:
				ind.SetFloat(0)
			}
		}
	}
}

func toString(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}
