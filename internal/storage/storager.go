package storage

import "github.com/thefrol/kysh-kysh-meow/internal/metrica"

// DEPRECATED todo delete
//
// Storager это интерфейс хранилища данных для метрик,
// Каждый новый тип метрики должен добавлять свой интерфейс сюда
type Storager interface {
	SetGauge(name string, v metrica.Gauge)
	Gauge(name string) (metrica.Gauge, bool)
	ListGauges() []string

	SetCounter(name string, v metrica.Counter)
	Counter(name string) (metrica.Counter, bool)
	ListCounters() []string

	Metricas() []metrica.Metrica
}
