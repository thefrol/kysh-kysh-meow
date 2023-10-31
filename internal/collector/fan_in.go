package collector

import (
	"context"
	"sync"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

func FanIn(ctx context.Context, chs ...<-chan metrica.Metrica) chan metrica.Metrica {
	chMix := make(chan metrica.Metrica, len(chs)*generatorChannelSize)
	fanWG := sync.WaitGroup{} // исплючительно внутренняя группа ожидания

	// cоздаем воркера под каждый канал
	for _, ch := range chs {

		fanWG.Add(1)

		go func(ch <-chan metrica.Metrica) {
			defer fanWG.Done()

			for v := range ch {
				select {
				case chMix <- v:

				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	// горутина будет ожидать закрытия канала
	wg.Add(1)
	go func() {
		fanWG.Wait()
		close(chMix)
		wg.Done() // todo потом эти все ожидания можно будет убрать, потому что по сути нам надо
		// дождаться только остановки воркеров
	}()

	return chMix
}
