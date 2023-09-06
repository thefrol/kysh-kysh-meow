package metrica

import "fmt"

type (
	Gauge float64
)

func (g Gauge) String() string {
	return fmt.Sprint(float64(g))
}

// Проверка, что метрика соответсвует нужному интерфейсу
var _ fmt.Stringer = (*Gauge)(nil)
