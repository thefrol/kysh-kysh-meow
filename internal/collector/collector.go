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
	"github.com/thefrol/kysh-kysh-meow/lib/scheduler"
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

	// запуск планировщика
	c := scheduler.New()

	// отправляем данные раз в repostInterval
	c.AddJob(time.Duration(config.ReportInterval)*time.Second, func() {
		//отправляем на сервер метрики из хранилища s
		batch := []metrica.Metrica{}
		log.Debug().Msg("Читаю метрики из канала")
	readLoop:
		for {
			select {
			case m := <-inMs:
				batch = append(batch, m)
			case m := <-inPc:
				batch = append(batch, m)
			case m := <-inRv:
				batch = append(batch, m)
			case m := <-inPs:
				batch = append(batch, m)
			default:
				break readLoop
			}
		}
		// отправляем
		log.Debug().Int("batch_len", len(batch)).Msg("Отправляю метрики")
		err := report.Send(batch, Endpoint(config.Addr, updateRoute))
		if err != nil {
			log.Error().Msgf("Попытка отправить метрики завершилась с  ошибками: %v", err)
			return
		}

	})

	// Запускаем планировщик, и он занимает поток
	c.Serve(200 * time.Millisecond)

	wg.Wait()
}

// Endpoint формирует точку, куда агент будет посылать все запросы на основе своей текущей конфигурации
func Endpoint(addr, route string) string {
	return fmt.Sprintf("%s%s", "http://", path.Join(addr, route))
}
