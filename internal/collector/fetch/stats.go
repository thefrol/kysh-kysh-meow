package fetch

import (
	"runtime"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// Stats в терминологии DDD представляет структуру данных, которую можно будет преобразовить
// в энтити metrica.Metrica.
//
// это такой сырой, полуоформленный формат данных. Я не разбираю его сразу,
// просто потому что не каждый опрос памяти будет отправлен. Не охота
// тратить на это время и оперативную память.
type Stats struct {
	memStats    *runtime.MemStats
	randomValue metrica.Gauge
	pollCount   metrica.Metrica
}

// Преобразует хранящиеся значения в транспортную структуру metrica.Metrica
func (st Stats) ToTransport() (m []metrica.Metrica) {
	m = append(m, metrica.Gauge(st.memStats.Alloc).Metrica("Alloc"))
	m = append(m, metrica.Gauge(st.memStats.BuckHashSys).Metrica("BuckHashSys"))
	m = append(m, metrica.Gauge(st.memStats.Frees).Metrica("Frees"))
	m = append(m, metrica.Gauge(st.memStats.GCCPUFraction).Metrica("GCCPUFraction"))
	m = append(m, metrica.Gauge(st.memStats.GCSys).Metrica("GCSys"))
	m = append(m, metrica.Gauge(st.memStats.HeapAlloc).Metrica("HeapAlloc"))
	m = append(m, metrica.Gauge(st.memStats.HeapIdle).Metrica("HeapIdle"))
	m = append(m, metrica.Gauge(st.memStats.HeapInuse).Metrica("HeapInuse"))
	m = append(m, metrica.Gauge(st.memStats.HeapObjects).Metrica("HeapObjects"))
	m = append(m, metrica.Gauge(st.memStats.HeapReleased).Metrica("HeapReleased"))
	m = append(m, metrica.Gauge(st.memStats.HeapSys).Metrica("HeapSys"))
	m = append(m, metrica.Gauge(st.memStats.LastGC).Metrica("LastGC"))
	m = append(m, metrica.Gauge(st.memStats.Lookups).Metrica("Lookups"))
	m = append(m, metrica.Gauge(st.memStats.MCacheInuse).Metrica("MCacheInuse"))
	m = append(m, metrica.Gauge(st.memStats.MCacheSys).Metrica("MCacheSys"))
	m = append(m, metrica.Gauge(st.memStats.MSpanInuse).Metrica("MSpanInuse"))
	m = append(m, metrica.Gauge(st.memStats.MSpanSys).Metrica("MSpanSys"))
	m = append(m, metrica.Gauge(st.memStats.Mallocs).Metrica("Mallocs"))
	m = append(m, metrica.Gauge(st.memStats.NextGC).Metrica("NextGC"))
	m = append(m, metrica.Gauge(st.memStats.NumForcedGC).Metrica("NumForcedGC"))
	m = append(m, metrica.Gauge(st.memStats.NumGC).Metrica("NumGC"))
	m = append(m, metrica.Gauge(st.memStats.OtherSys).Metrica("OtherSys"))
	m = append(m, metrica.Gauge(st.memStats.PauseTotalNs).Metrica("PauseTotalNs"))
	m = append(m, metrica.Gauge(st.memStats.StackInuse).Metrica("StackInuse"))
	m = append(m, metrica.Gauge(st.memStats.StackSys).Metrica("StackSys"))
	m = append(m, metrica.Gauge(st.memStats.Sys).Metrica("Sys"))
	m = append(m, metrica.Gauge(st.memStats.TotalAlloc).Metrica("TotalAlloc"))

	// случайное значение
	m = append(m, st.randomValue.Metrica(randomValueName))

	// счетчик
	m = append(m, st.pollCount)

	return
}
