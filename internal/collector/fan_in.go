package collector

import (
	"sync"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// FanIn возвращает канал, в который будут сливаться данные из указанных
// каналов chs. Когда все входные каналы будут закрыты, то и выходной тоже закроется.
func FanIn(chs ...<-chan metrica.Metrica) chan metrica.Metrica {
	chMix := make(chan metrica.Metrica, len(chs)*generatorChannelSize)
	wg := sync.WaitGroup{}

	// cоздаем воркера под каждый канал
	for _, ch := range chs {

		wg.Add(1)

		go func(ch <-chan metrica.Metrica) {
			defer wg.Done()

			for v := range ch {
				chMix <- v // что будет если case a<-channel, в каких случаях пропускает?
			}
		}(ch)
	}

	// горутина будет ожидать закрытия канала
	go func() {
		wg.Wait()
		close(chMix)
	}()

	return chMix
}
