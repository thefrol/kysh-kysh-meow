package collector

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/fetch"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

type FetchFunc func() fetch.Batcher

// generator это источник данных. Раз в timeout секунд он опрашивает источник fetch,
// формирует данные и отправляет в исходящий канал. Возвращает исходящий канал
// когда истечет контекст ctx, канал будет закрыт.
func generator(ctx context.Context, fetch FetchFunc, timeout time.Duration) <-chan metrica.Metrica {
	chGen := make(chan metrica.Metrica, generatorChannelSize)
	tick := time.NewTicker(timeout)

	go func() {
		for {
			select {
			case <-tick.C:
				// собираем метрики

				log.Debug().Msg("Опрашиваются метрики") // сюда бы имя добавить какое)

				// я вдруг подумал, что логгирование
				// ещё неплохо документирует сам код

				ms := fetch().ToTransport()
				for _, m := range ms {
					// todo. а что будет если канал вдруг переполнен?
					chGen <- m
				}
			case <-ctx.Done():
				close(chGen)
				tick.Stop()
				return
			}

		}
	}()
	return chGen
}
