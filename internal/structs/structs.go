// structs содержит инструменты рефлекции для работы со структурами
// 	FieldsFloat: возвращает мапу из полей, которые можно преобразовать во float64
// 	FieldNames: возвращает коллекцию имен полей
package structs

import (
	"errors"
	"reflect"
)

var (
	ErrorNotStruct  = errors.New("received object is not a struct")
	ErrorNilPointer = errors.New("received nil pointer")
)

// getFieldsFloat возвращает из структуры s поля, приведенные к указанному типу.
// Возвращены будут только поля с именами, неуказанными в exclude
// Если exclude=nil, то фильтрации не производится
// Возвращает ошибку в случае ошибки
func FieldsFloat(s interface{}, exclude []string) (m map[string]float64, err error) {
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
		switch v := r.Field(i).Interface().(type) { // if is convertible #TODO
		case int64:
			m[r.Type().Field(i).Name] = float64(v)
		case int32:
			m[r.Type().Field(i).Name] = float64(v)
		case uint64:
			m[r.Type().Field(i).Name] = float64(v)
		case uint32:
			m[r.Type().Field(i).Name] = float64(v)
		case float64:
			m[r.Type().Field(i).Name] = v
		case float32:
			m[r.Type().Field(i).Name] = float64(v)
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

func FieldNames(s interface{}) (names []string, err error) {
	r := reflect.ValueOf(s)
	if reflect.Indirect(r).Kind() != reflect.Struct { // будет что потестить // make a restricting interface!
		return nil, ErrorNotStruct
	}
	for i := 0; i < r.NumField(); i++ {
		names = append(names, r.Type().Field(i).Name)
	}
	return names, nil
}
