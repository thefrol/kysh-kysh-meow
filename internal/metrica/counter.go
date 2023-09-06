package metrica

import "fmt"

type (
	Counter int64
)

func (c Counter) String() string {
	return fmt.Sprint(int64(c))
}

// Проверка, что метрика соответсвует нужному интерфейсу
var _ fmt.Stringer = (*Counter)(nil)
