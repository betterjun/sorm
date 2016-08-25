package sorm

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"
)

type result struct {
	rows *sql.Rows
	cols []string
}

func (r *result) Next(obj interface{}, args ...interface{}) (err error) {
	if r.rows == nil {
		return fmt.Errorf("result is not initialized")
	}
	if r.cols == nil {
		r.cols, err = r.rows.Columns()
		if err != nil {
			return err
		}
	}

	if r.rows.Next() {
		scanArgs, err := getFieldsForOne(obj, args, r.cols)
		if err != nil {
			return err
		}
		if scanArgs == nil {
			return fmt.Errorf("no receiver found")
		}

		return r.rows.Scan(scanArgs...)
	} else {
		r.rows.Close()
		return io.EOF
	}
}

func (r *result) All(objs interface{}) (err error) {
	if r.rows == nil {
		return fmt.Errorf("result is not initialized")
	}

	val := reflect.ValueOf(objs)
	sInd := reflect.Indirect(val)
	if val.Kind() != reflect.Ptr || sInd.Kind() != reflect.Slice {
		return fmt.Errorf("receiver must be a pointer of slice")
	}

	cols, err := r.rows.Columns()
	if err != nil {
		return err
	}
	defer r.rows.Close()

	var scanArgs []interface{}
	var ind reflect.Value
	etyp := sInd.Type().Elem()
	if etyp.Kind() == reflect.Struct {
		ind = reflect.New(sInd.Type().Elem()).Elem()
		scanArgs = getScanFieldFromStruct(ind, cols)
		if scanArgs == nil {
			return fmt.Errorf("no receiver found")
		}
	} else {
		if len(cols) > 1 {
			return fmt.Errorf("the result has more than one coloum, please passing in a struct slice")
		}

		ind = reflect.New(etyp).Elem()
		scanArgs = []interface{}{ind.Addr().Interface()}
	}

	sIndCopy := sInd
	err = io.EOF
	for r.rows.Next() {
		err = r.rows.Scan(scanArgs...)
		if err != nil {
			break
		}

		sIndCopy = reflect.Append(sIndCopy, ind)
	}

	// ret may take back some records, even though there is an error.
	sInd.Set(sIndCopy)
	return err
}

func (r *result) ColumnNames() (cols []string, err error) {
	if r.rows == nil {
		return nil, fmt.Errorf("result is not initialized")
	}

	if r.cols == nil {
		r.cols, err = r.rows.Columns()
		if err != nil {
			return nil, err
		}
	}
	cols = make([]string, len(r.cols))
	copy(cols, r.cols)
	return cols, nil
}

func (r *result) Close() (err error) {
	if r.rows != nil {
		err = r.rows.Close()
		r.rows = nil
	}
	r.cols = nil
	return err
}
