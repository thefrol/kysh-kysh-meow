package stats

import (
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

var pollCount metrica.Counter

// DropPoll сбрасывает значение счетчика опросов памяти в указанном хранилище
// dropCounter сбраcывает счетчик
func DropPollCount() {
	pollCount = 0
}
func incrementPollCount() {
	pollCount += metrica.Counter(1)
}
