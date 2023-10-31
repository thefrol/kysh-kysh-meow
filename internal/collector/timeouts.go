package collector

import (
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// WithTimeouts возвращает outCh, в который будут поступать данные из inCh,
// только с интервалами. Если inCh пустеет, то мы ждем не менее timeout секунд
// и продолжаем работу
func WithTimeouts(inCh <-chan metrica.Metrica, timeout time.Duration) chan metrica.Metrica {
	chOut := make(chan metrica.Metrica, MaxBatch)

	tick := time.NewTicker(timeout)

	go func() {
	loop:
		for {
			select {
			case v, ok := <-inCh:
				if !ok {
					// если канал закрылся, то выйдем
					// из цикла и передем к останове
					break loop
				}
				chOut <- v
			default:
				// если все данные прочитали, то ждем следующего таймера
				<-tick.C
			}
		}
		// входной канал закрылся,
		// закроем и мы свой

		close(chOut)

		// думаю, тут можно было бы и написать как-то через range tick.C
		// но канал тикера как бы не закрывается никогда(((
		//
		// но мы все равно делаем выход через if

	}()
	return chOut
}
