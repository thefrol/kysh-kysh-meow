package storage

type MemStore struct {
	counters map[string]Counter
	gauges   map[string]Gauge
}

// Возвращает значание именнованного счетчика, и булево значение: удалось ли счетчик найти
func (m MemStore) Counter(name string) (Counter, bool) {
	val, ok := m.counters[name] //Можно ли сделать в одну строчку #MENTOR
	return val, ok
}

func (m MemStore) SetCounter(name string, value Counter) {
	m.counters[name] = value
}

func (m MemStore) Gauge(name string) (Gauge, bool) {
	val, ok := m.gauges[name] //Можно ли сделать в одну строчку #MENTOR
	return val, ok
}

func (m MemStore) SetGauge(name string, value Gauge) {
	m.gauges[name] = value
}

// Проверка, что MemStore соответсвует нужному интерфейсу
var _ Storager = (*MemStore)(nil)
