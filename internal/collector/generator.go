package collector

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/fetch"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

type FetchFunc func() fetch.Batcher

var wg = sync.WaitGroup{}

func generator(ctx context.Context, fetch FetchFunc, timeout time.Duration) <-chan metrica.Metrica {
	chGen := make(chan metrica.Metrica, generatorChannelSize)

	wg.Add(1)

	tick := time.NewTicker(timeout)

	go func() {

		for {
			select {
			case <-tick.C:
				// собираем метрики
				// не хорошо будет тут зависнуть, конечно. Нужно чтобы
				// отправщики последними останавливались
				log.Debug().Msg("Собираются метрики") // сюда бы имя добавить какое)

				ms := fetch().ToTransport()
				for _, m := range ms {
					chGen <- m
				}
			case <-ctx.Done():
				close(chGen) // кто создал тот и закрывает
				wg.Done()
			}

		}
	}()
	return chGen
}
