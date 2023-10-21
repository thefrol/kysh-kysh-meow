// slices помогает работать со слайсами и содержит вспомогательные функции
// 	contains: проверяет на вхождение элемента в слайс
// 	containsSlice: проверяет на вхождение целого слайса в другой слайс
package slices

// Contains созвращает True, если переданный слайс содержит элемент value
func Contains[T comparable](s []T, value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

// ContainsSlice возвращет true, если слайс a содержит в себя слайс b
func СontainsSlice[T comparable](sl []T, b []T) bool {
	if len(b) > len(sl) {
		return false // b is bigger
	}
	for _, vb := range b {
		if !Contains[T](sl, vb) {
			return false
		}
	}
	return true
}

// Difference убирает из слайса элекменты другого слайса
func Difference[T comparable](from []T, exclude []T) (diff []T) {
	for _, v := range from { //как это все поэлегантней то сделать
		if !Contains[T](exclude, v) {
			diff = append(diff, v)
		}
	}
	return
}
