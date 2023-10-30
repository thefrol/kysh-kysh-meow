package fetch

import (
	"github.com/thefrol/kysh-kysh-meow/internal/collector/internal/pollcount"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/internal/randomvalue"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

func PollCount() Batcher {
	pollcount.Drop() // сбрасываем каждый как раз прочитали, ну что-то тут происходит короче
	return Single(pollcount.Get())
}

func RandomValue() Batcher {
	return Single(randomvalue.Get())
}

// Single это контейнер для одной единственной метрики, но так, чтобы
// удовлетворяла batcher
type Single metrica.Metrica

func (s Single) ToTransport() []metrica.Metrica {
	return []metrica.Metrica{metrica.Metrica(s)}
}
