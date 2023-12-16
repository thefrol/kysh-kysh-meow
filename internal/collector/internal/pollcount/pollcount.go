// Этот пакет содержит логику работы счетчика PollCount в агенте
// Он специально лежит в internal, чтобы нельзя было поменять из агента
// значение счетчика
package pollcount

import (
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// todo нужен atomic
var pollCount int64

const IDPollCount = "PollCount"

// Drop сбрасывает значение счетчика опросов памяти
func Drop() {
	pollCount = 0
}

func Increment() {
	pollCount += 1
}

func Get() metrica.Metrica {
	val := pollCount
	return metrica.Metrica{
		MType: metrica.CounterName,
		ID:    IDPollCount,
		Delta: &val,
	}
}
