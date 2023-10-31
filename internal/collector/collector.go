package collector

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/fetch"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/report"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const (
	generatorChannelSize = 300
)

func FetchAndReport(config config.Agent, updateRoute string) {

	// КОРОЧЕ
	//
	// будет так
	//
	// Тра-ля-ля поллкаунт там как-то обноваляется. для него отдельная горутина
	// Рендом велью тоже. Они отправляются по одному как бы
	//
	// Для мемстатс третья горутина
	//
	// И для новых метрик четвертая

	// Должен быть какой-то метод типа newGenerator, который создает такие потоки

	// Эти значит генераторы пишут в канал свой дорогой метрики по одной,
	// а воркеры их собирают в пачки и отправляют
	// а может пачки для них кто-то другой подготавливает даже

	// Сделать классную такую мермаид диаграмму со всеми каналами, кто куда как собирает

	// создать каналы сбора метрик
	interval := time.Second * time.Duration(config.PollingInterval)
	ctx := context.TODO()
	inMs := generator(ctx, fetch.MemStats, interval)
	inPc := generator(ctx, fetch.PollCount, interval)
	inRv := generator(ctx, fetch.RandomValue, interval)
	inPs := generator(ctx, fetch.GoPS, interval)

	// объединить каналы в один
	inMix := FanIn(ctx, inMs, inPc, inPs, inRv)

	// запуск планировщика
	reportInterval := time.Second * time.Duration(config.ReportInterval)
	pool(ctx, 2, func() {
		batch := []metrica.Metrica{}
		log.Debug().Msg("Читаю метрики из канала")

		tick := time.NewTicker(reportInterval)

	readLoop:
		for {
			select {
			case m := <-inMix:
				batch = append(batch, m)
				if len(batch) >= MaxBatch {
					break readLoop
				}
			case <-tick.C:
				// отправляем

				// todo по хорошему именно вот эту штуку бы в горутине обрабатывать, пусть именно она в пуле висит
				// а мы раз в десять секунд собираем в батчи и передаем дальше.
				// типа добавим ещё один элемент конвеера
				log.Debug().Int("batch_len", len(batch)).Msg("Отправляю метрики")
				err := report.Send(batch, Endpoint(config.Addr, updateRoute))
				if err != nil {
					log.Error().Msgf("Попытка отправить метрики завершилась с  ошибками: %v", err)
					return
				}
			}
		}
	})

	wg.Wait()
}

// Endpoint формирует точку, куда агент будет посылать все запросы на основе своей текущей конфигурации
func Endpoint(addr, route string) string {
	return fmt.Sprintf("%s%s", "http://", path.Join(addr, route))
}
