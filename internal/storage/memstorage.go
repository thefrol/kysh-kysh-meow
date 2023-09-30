package storage

import "github.com/thefrol/kysh-kysh-meow/internal/metrica"

type MemStore struct {
	Counters map[string]metrica.Counter
	Gauges   map[string]metrica.Gauge
}

func New() (m MemStore) {
	m.Counters = make(map[string]metrica.Counter)
	m.Gauges = make(map[string]metrica.Gauge)
	return
}

// Возвращает значание именнованного счетчика, и булево значение: удалось ли счетчик найти
func (m MemStore) Counter(name string) (metrica.Counter, bool) {
	val, ok := m.Counters[name] //Можно ли сделать в одну строчку #MENTOR
	return val, ok
}

func (m MemStore) SetCounter(name string, value metrica.Counter) {
	m.Counters[name] = value
}

func (m MemStore) ListCounters() (keys []string) {
	//keys = make([]string, 0, len(m.counters)) не такая плохая идея
	for k := range m.Counters {
		keys = append(keys, k)
	}
	return
}

func (m MemStore) Gauge(name string) (metrica.Gauge, bool) {
	val, ok := m.Gauges[name] //Можно ли сделать в одну строчку #MENTOR
	return val, ok
}

func (m MemStore) SetGauge(name string, value metrica.Gauge) {
	m.Gauges[name] = value
}

func (m MemStore) ListGauges() (keys []string) {
	//keys = make([]string, 0, len(m.counters)) не такая плохая идея
	for k := range m.Gauges {
		keys = append(keys, k)
	}
	return
}

func (m MemStore) Metricas() (list []metrica.Metrica) {
	for k, c := range m.Counters {
		list = append(list, c.Metrica(k))
	}
	for k, g := range m.Gauges {
		list = append(list, g.Metrica(k))
	}
	return
}

// todo
//
// У меня может быть стринг интерфейс к хранилищу, и им пользуются всякие ребятки, но другие ребятки им не пользуются. Но дело ли это хранилища,
// Или можно какую-то надсткойку сделать, например, - новый слой?
// Проверка, что MemStore соответсвует нужному интерфейсу
var _ Storager = (*MemStore)(nil)
