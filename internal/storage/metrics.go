package storage

import "fmt"

type (
	Counter int64
	Gauge   float64
)

func (c Counter) String() string {
	return fmt.Sprint(int64(c))
}

func (g Gauge) String() string {
	return fmt.Sprint(float64(g))
}

// Проверка, что Метрики соответсвует нужному интерфейсу
var _ fmt.Stringer = (*Counter)(nil)
var _ fmt.Stringer = (*Gauge)(nil)
