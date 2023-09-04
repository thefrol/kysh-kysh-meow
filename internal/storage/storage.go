package storage

type MemStore struct {
	counters map[string]Counter
	gauges   map[string]Gauge
}

func New() (m MemStore) {
	m.counters = make(map[string]Counter)
	m.gauges = make(map[string]Gauge)
	return
}

// Возвращает значание именнованного счетчика, и булево значение: удалось ли счетчик найти
func (m MemStore) Counter(name string) (Counter, bool) {
	val, ok := m.counters[name] //Можно ли сделать в одну строчку #MENTOR
	return val, ok
}

func (m MemStore) SetCounter(name string, value Counter) {
	m.counters[name] = value
}

func (m MemStore) ListCounters() (keys []string) {
	//keys = make([]string, 0, len(m.counters)) не такая плохая идея
	for k := range m.counters {
		keys = append(keys, k)
	}
	return
}

func (m MemStore) Gauge(name string) (Gauge, bool) {
	val, ok := m.gauges[name] //Можно ли сделать в одну строчку #MENTOR
	return val, ok
}

func (m MemStore) SetGauge(name string, value Gauge) {
	m.gauges[name] = value
}

func (m MemStore) ListGauges() (keys []string) {
	//keys = make([]string, 0, len(m.counters)) не такая плохая идея
	for k := range m.gauges {
		keys = append(keys, k)
	}
	return
}

// Проверка, что MemStore соответсвует нужному интерфейсу
var _ Storager = (*MemStore)(nil)
