package collector

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// TimeGate возвращает outCh, в который будут переброшены данные из inCh.
// Если inCh пустеет, то мы засыпаем на не более timeout секунд, и потом
// продолжаем переброску из одного канала в другой
func TimeGate(ctx context.Context, inCh <-chan metrica.Metrica, timeout time.Duration) chan metrica.Metrica {
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
				// Eсли все данные прочитали, то ждем следующего таймера
				// или если завершился контекст,
				// разрешаем быстро пропустить
				select {
				case <-tick.C:
				case <-ctx.Done():
				}
			}
		}

		// входной канал закрылся,
		// закроем и мы свой
		// и подчистим концы
		tick.Stop()

		close(chOut)
		log.Debug().Msg("WithTimeout завершил работу")

	}()
	return chOut
}
