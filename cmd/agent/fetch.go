// fetchStats собирает метрики из памяти, и выдает их удобной мапой
package main

import (
	"errors"
	"reflect"
)

func fetchStats() (m map[string]float64) {
	return
}

var (
	ErrorNotStruct  = errors.New("received object is not a struct")
	ErrorNilPointer = errors.New("received nil pointer")
)

func getFieldsFloat(s interface{}) (m map[string]float64, err error) {
	r := reflect.ValueOf(s)
	if reflect.Indirect(r).Kind() != reflect.Struct { // будет что потестить // make a restricting interface!
		return nil, ErrorNotStruct
	}
	// check if pointer
	// if r.IsNil() {
	// 	return nil, ErrorNilPointer

	// }

	m = make(map[string]float64)
	for i := 0; i < r.NumField(); i++ {
		switch v := r.Field(i).Interface().(type) { // if is convertible
		case int64:
			m[r.Type().Field(i).Name] = float64(v)
		case uint64:
			m[r.Type().Field(i).Name] = float64(v)
		case float64:
			m[r.Type().Field(i).Name] = v
		default:
			continue
		}
	}

	//check if assignable! most of values are int64
	return
}
