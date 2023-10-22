// Этот пакет содержит логику работы счетчика PollCount в агенте
// Он специально лежит в internal, чтобы нельзя было поменять из агента
// значение счетчика
package pollcount

var pollCount int64

// Drop сбрасывает значение счетчика опросов памяти
func Drop() {
	pollCount = 0
}

func Increment() {
	pollCount += 1
}

func Get() int64 {
	return pollCount
}
