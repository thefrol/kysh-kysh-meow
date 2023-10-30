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
	"math/rand"
	"runtime"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/collector/internal/pollcount"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const (
	// название рамндомной метрики среди всех данных, что мы собираем
	randomValueName = "RandomValue"
)

// MemStats собирает метрики мамяти и сохраняет их во временное хранилище
func MemStats() Stats {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	s := Stats{
		memStats:    &m,
		randomValue: randomGauge(),
		pollCount:   pollcount.Get(),
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
