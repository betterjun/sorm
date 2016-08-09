package sorm

import (
	"fmt"
	"reflect"
	"strings"
)

func getScanFields(ptr interface{}, cols []string) (scanArgs []interface{}) {
	fields := make(map[string]interface{})
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("orm")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		// take the addr and as interface{}, this is required by sql Scan
		fields[name] = v.Field(i).Addr().Interface()
	}

	for _, name := range cols {
		f := fields[name]
		scanArgs = append(scanArgs, f)
	}
	fmt.Println(scanArgs)
	return scanArgs
}
