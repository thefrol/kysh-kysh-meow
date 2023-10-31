package collector

import (
	"context"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

func FanIn(ctx context.Context, chs ...<-chan metrica.Metrica) chan metrica.Metrica {
	chMix := make(chan metrica.Metrica, len(chs)*generatorChannelSize)

	// cоздаем воркера под каждый канал
	for _, ch := range chs {
		defer wg.Done()
		wg.Add(1)

		go func(ch <-chan metrica.Metrica) {

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
	go func() {
		wg.Wait()
		close(chMix)
	}()

	return chMix
}
