package storage

// Storager это интерфейс хранилища данных для метрик,
// Каждый новый тип метрики должен добавлять свой интерфейс сюда
type Storager interface {
	Gauger
	Counterer
}

type Gauger interface {
	SetGauge(name string, v Gauge)
	Gauge(name string) (Gauge, bool)
	ListGauges() []string
}

type Counterer interface {
	SetCounter(name string, v Counter)
	Counter(name string) (Counter, bool)
	ListCounters() []string
}
