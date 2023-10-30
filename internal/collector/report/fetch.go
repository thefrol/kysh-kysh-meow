package report

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/collector/internal/pollcount"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const (
	// название рамндомной метрики среди всех данных, что мы собираем
	randomValueName = "RandomValue"
	// Счетчик поличества опросов
	metricPollCount = "PollCount"
)

// Fetch собирает метрики мамяти и сохраняет их во временное хранилище
func Fetch() Stats {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	s := Stats{
		memStats:    &m,
		randomValue: randomGauge(),
		pollCount:   metrica.Counter(pollcount.Get()),
	}

	// Добавить ко счетчику опросов
	pollcount.Increment()

	return s
}

// randomGauge возвращает случайное число типа float64
func randomGauge() metrica.Gauge {
	s := rand.NewSource(int64(time.Now().Nanosecond()))
	r := rand.New(s)
	return metrica.Gauge(r.Float64())
}
