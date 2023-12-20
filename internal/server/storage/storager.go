package storage

import "github.com/thefrol/kysh-kysh-meow/internal/metrica"

// legacyStorager это интерфейс хранилища данных для метрик,
// используется только для совместимости со старыми видами хранилищ,
// использвуется чтобы оборачивать с помощью NewAdapter()
type legacyStorager interface {
	SetGauge(name string, v metrica.Gauge)
	Gauge(name string) (metrica.Gauge, bool)
	ListGauges() []string

	SetCounter(name string, v metrica.Counter)
	Counter(name string) (metrica.Counter, bool)
	ListCounters() []string

	Metricas() []metrica.Metrica
}
