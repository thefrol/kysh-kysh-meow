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
	"github.com/thefrol/kysh-kysh-meow/internal/collector/internal/randomvalue"
)

// MemStats собирает метрики мамяти и сохраняет их во временное хранилище
func MemStats() Batcher {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	s := Stats{
		memStats:    &m,
		randomValue: randomvalue.Get(),
		pollCount:   pollcount.Get(),
	}

	// Добавить ко счетчику опросов
	pollcount.Increment()

	return s
}
