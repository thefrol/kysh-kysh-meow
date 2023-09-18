// stats собирает метрики из памяти, и выдает их удобной мапой
package stats

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

const randomValueName = "RandomValue"

// fetchMemStats собирает метрики мамяти и сохраняет их в хранилище store
func FetchMemStats(s storage.Storager) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	s.SetGauge("Alloc", metrica.Gauge(m.Alloc))
	s.SetGauge("BuckHashSys", metrica.Gauge(m.BuckHashSys))
	s.SetGauge("Frees", metrica.Gauge(m.Frees))
	s.SetGauge("GCCPUFraction", metrica.Gauge(m.GCCPUFraction))
	s.SetGauge("GCSys", metrica.Gauge(m.GCSys))
	s.SetGauge("HeapAlloc", metrica.Gauge(m.HeapAlloc))
	s.SetGauge("HeapIdle", metrica.Gauge(m.HeapIdle))
	s.SetGauge("HeapInuse", metrica.Gauge(m.HeapInuse))
	s.SetGauge("HeapObjects", metrica.Gauge(m.HeapObjects))
	s.SetGauge("HeapSys", metrica.Gauge(m.HeapSys))
	s.SetGauge("LastGC", metrica.Gauge(m.LastGC))
	s.SetGauge("Lookups", metrica.Gauge(m.Lookups))
	s.SetGauge("MCacheInuse", metrica.Gauge(m.MCacheInuse))
	s.SetGauge("MCacheSys", metrica.Gauge(m.MCacheSys))
	s.SetGauge("MSpanInuse", metrica.Gauge(m.MSpanInuse))
	s.SetGauge("MSpanSys", metrica.Gauge(m.MSpanSys))
	s.SetGauge("Mallocs", metrica.Gauge(m.Mallocs))
	s.SetGauge("NextGC", metrica.Gauge(m.NextGC))
	s.SetGauge("NumForcedGC", metrica.Gauge(m.NumForcedGC))
	s.SetGauge("NumGC", metrica.Gauge(m.NumGC))
	s.SetGauge("OtherSys", metrica.Gauge(m.OtherSys))
	s.SetGauge("PauseTotalNs", metrica.Gauge(m.PauseTotalNs))
	s.SetGauge("StackInuse", metrica.Gauge(m.StackInuse))
	s.SetGauge("StackSys", metrica.Gauge(m.StackSys))
	s.SetGauge("Sys", metrica.Gauge(m.Sys))
	s.SetGauge("TotalAlloc", metrica.Gauge(m.TotalAlloc))

	// случайное значение
	s.SetGauge(randomValueName, metrica.Gauge(randomFloat64()))

	// тут, конечно тоже можно наошибаться, и может рефлексия поможет для тестов
	// но тут уже можно взять сторонний набор
}

// randomFloat64 возвращает случайное число типа float64
func randomFloat64() float64 {
	s := rand.NewSource(int64(time.Now().Nanosecond()))
	r := rand.New(s)
	return r.Float64()
}
