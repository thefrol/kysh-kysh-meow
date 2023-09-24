package metrica

import "fmt"

type Counter int64

func (c Counter) String() string {
	return fmt.Sprint(int64(c))
}

func (c Counter) Metrica(id string) Metrica {
	val := int64(c) // todo дпоменять бы типа для Metrica
	return Metrica{
		MType: CounterName,
		ID:    id,
		Delta: &val,
	}
}

func (c *Counter) ParseMetrica(m Metrica) {
	*c = Counter(*m.Delta)
}

// Проверка, что метрика соответсвует нужному интерфейсу
var _ fmt.Stringer = (*Counter)(nil)
var _ Metrer = (*Counter)(nil)
