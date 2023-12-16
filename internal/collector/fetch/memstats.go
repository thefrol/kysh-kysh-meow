// Этот пакет отвечает за получение метрик из других пакетов операционной
// системы и трансформацию их с понятный для нашего сервиса формат Metrica
//
// Значение большинства метрик забираеются из пакета runtime.ReadMemStats() и ещё дополнительно пополняются двумя параметрами,
// один из которых счетчик PollCount - там хранится число раз, сколько мы опросили память. После отправки на сервер данных, это
// значение сбрасывается
//
// Основые функции:
//
// Memstats() - собирает основые параметры использования памяти и сохраняет
// в промеждуточное хранилише типа Stats, а так же увеличивает счетчик PollCount
package fetch

import (
	"runtime"

	"github.com/thefrol/kysh-kysh-meow/internal/collector/internal/pollcount"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// MemStats собирает метрики мамяти и сохраняет их во временное хранилище
func MemStats() Batcher {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	s := MemBatch{
		memStats: &m,
	}

	// Добавить ко счетчику опросов
	pollcount.Increment()

	return s
}

// MemBatch в терминологии DDD представляет структуру данных, которую можно будет преобразовить
// в энтити metrica.Metrica.
//
// это такой сырой, полуоформленный формат данных. Я не разбираю его сразу,
// просто потому что не каждый опрос памяти будет отправлен. Не охота
// тратить на это время и оперативную память.
type MemBatch struct {
	memStats *runtime.MemStats
}

// Преобразует хранящиеся значения в транспортную структуру metrica.Metrica
func (st MemBatch) ToTransport() (m []metrica.Metrica) {
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
	return
}
