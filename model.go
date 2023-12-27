package kyshkyshmeow

import (
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/metricas"
)

// Gauge используется чтобы передавать
// изменение величины
type Gauge struct {
	// имя величины
	ID string

	// новое значение
	Value float64
}

func (g Gauge) toM() metricas.Metrica {
	return metrica.Metrica{
		MType: metrica.GaugeName,
		ID:    g.ID,
		Value: &g.Value, // [bug] отправится в кучу
	}
}

// Counter используется чтобы передавать
// изменение счетчика
type Counter struct {
	// имя счетчика
	ID string

	// изменение счетчика
	Delta int64
}

func (c Counter) toM() metricas.Metrica {
	// BUG(frolenkodima): при упаковке джейсон происходит много аллоков памяти
	return metrica.Metrica{
		MType: metrica.CounterName,
		ID:    c.ID,
		Delta: &c.Delta, // [bug] отправится в кучу
	}
}

// metrer представляет интерфейс для передаваемой на сервер метрики
type metrer interface {
	toM() metrica.Metrica
}
