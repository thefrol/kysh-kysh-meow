// fetchStats собирает метрики из памяти, и выдает их удобной мапой
package main

import (
	"errors"
	"reflect"
	"runtime"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// fetchStats собирает метрики мамяти и сохраняет их в хранилище, исключая ненужные exclude
func saveMemStats(store storage.Storager, exclude []string) error {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	stats, err := getFieldsFloat(m, exclude)
	if err != nil {
		return err
	}
	for key, value := range stats {
		store.SetGauge(key, storage.Gauge(value)) //#TODO SetGauges()
	}
	return nil
}

var (
	ErrorNotStruct  = errors.New("received object is not a struct")
	ErrorNilPointer = errors.New("received nil pointer")
)

// getFieldsFloat возвращает из структуры s поля, приведенные к указанному типу.
// Возвращены будут только поля с именами, неуказанными в exclude
// Если exclude=nil, то фильтрации не производится
// Возвращает ошибку в случае ошибки
func getFieldsFloat(s interface{}, exclude []string) (m map[string]float64, err error) {
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

	//exclude
	if len(exclude) == 0 {
		return
	}
	for _, f := range exclude {
		delete(m, f)
	}
	return m, nil
}

func getStructFields(s interface{}) (names []string, err error) {
	r := reflect.ValueOf(s)
	if reflect.Indirect(r).Kind() != reflect.Struct { // будет что потестить // make a restricting interface!
		return nil, ErrorNotStruct
	}
	for i := 0; i < r.NumField(); i++ {
		names = append(names, r.Type().Field(i).Name)
	}
	return names, nil
}

// Difference убирает из слайса элекменты другого слайса
func Difference[T comparable](from []T, exclude []T) []T {
	for i, v := range from { //как это все поэлегантней то сделать
		if contains[T](exclude, v) {
			remove[T](from, i)
		}
	}
	return from //#TODO check if changes original slice #test
}

// contains созвращает True, если переданный слайс содержит элемент value
func contains[T comparable](s []T, value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

// containsSlice
func containsSlice[T comparable](a []T, b []T) bool {
	if len(b) > len(a) {
		return false // b is bigger
	}
	for _, vb := range b {
		if !contains[T](a, vb) {
			return false
		}
	}
	return true
}

// remove убирает элемент из слайса под номером i
func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
