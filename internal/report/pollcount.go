package report

import (
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

var pollCount metrica.Counter

// DropPoll сбрасывает значение счетчика опросов памяти в указанном хранилище
// dropCounter сбраcывает счетчик
func dropPollCount() {
	pollCount = 0 // todo поллкаунт может стать вообще внутренней штукой модуля report если их объдинить, тогда агенту вообще не надо будет об этом париться
}
func incrementPollCount() {
	pollCount += metrica.Counter(1)
}
