package collector

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/fetch"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

type FetchFunc func() fetch.Batcher

func generator(ctx context.Context, fetch FetchFunc, timeout time.Duration) <-chan metrica.Metrica {
	chGen := make(chan metrica.Metrica, generatorChannelSize)
	tick := time.NewTicker(timeout)

	go func() {
		for {
			select {
			case <-tick.C:
				// собираем метрики
				// не хорошо будет тут зависнуть, конечно. Нужно чтобы
				// отправщики последними останавливались
				log.Debug().Msg("Опрашиваются метрики") // сюда бы имя добавить какое)

				// я вдруг подумал, что логгирование
				// ещё неплохо документирует сам код

				ms := fetch().ToTransport()
				for _, m := range ms {
					chGen <- m
				}
			case <-ctx.Done():
				close(chGen) // кто создал тот и закрывает
				tick.Stop()
				return
			}

		}
	}()
	return chGen
}
