package metrica

import "fmt"

type Gauge float64

func (g Gauge) String() string {
	return fmt.Sprint(float64(g))
}

func (g Gauge) Metrica(id string) Metrica {
	val := float64(g) // todo дпоменять бы типа для Metrica
	return Metrica{
		MType: "gauge",
		ID:    id,
		Value: &val,
	}
}

// Проверка, что метрика соответсвует нужным интерфейсам
var _ fmt.Stringer = (*Gauge)(nil)
var _ Metrer = (*Gauge)(nil)
