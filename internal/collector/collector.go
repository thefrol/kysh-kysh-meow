package collector

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/fetch"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/report"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const (
	generatorChannelSize = 300 // todo как вообще мне регулировать емкость этих каналов?
)

// FetchAndReport запускает основную логику. Источники данных, создаются воркеры,
// которые будут отправлять данные на url, из конфига cfg берутся данные по
// запуску(что вообще говоря некорректно todo, не хочу чтобы сервис зависел от конфига)
// И завершается когда истечет контекст ctx.
//
// Занимает текущую горутину, и освобождает её, когда закроется последний воркер.
// Это будет означать, что все созданные горутины штатно завершились и вся информация отправлена
//
// В текущей архитектуре мы гарантируем правильную отправку каунтеров, все каунтеры придут
// а очередность gauge может быть нарушена, возможно будет отправлено не самое свежее
// значение gauge в пределах retry.Timeout секунд
func FetchAndReport(ctx context.Context, config config.Agent, updateRoute string) {
	// создать каналы сбора метрик
	interval := time.Second * time.Duration(config.PollingInterval)
	inMs := generator(ctx, fetch.MemStats, interval)
	inPc := generator(ctx, fetch.PollCount, interval)
	inRv := generator(ctx, fetch.RandomValue, interval)
	inPs := generator(ctx, fetch.GoPS, interval)

	// объединить каналы в один
	inMix := FanIn(ctx, inMs, inPc, inPs, inRv)

	// собирать данные с перерывами
	reportInterval := time.Second * time.Duration(config.ReportInterval)
	inCh := TimeGate(ctx, inMix, reportInterval)

	// отправим данные
	workerCount := 3
	url := Endpoint(config.Addr, updateRoute)
	sema := NewSemaphore(int(config.RateLimit)) // регулирует максимальное количество исходяших соединений
	wg := sync.WaitGroup{}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(inCh, url, sema, &wg)
	}

	// todo
	//
	// надо конечно подумать над такими вещами, что если метрики не отправились? бросить их обратно в начало или в какую-то
	// Дополнительную очередь?
	// типа у нас ещё очеред в начало одна, мертвых сообщений они опять на входе

	wg.Wait()
}

// Endpoint формирует точку, куда агент будет посылать все запросы на основе своей текущей конфигурации
func Endpoint(addr, route string) string {
	return fmt.Sprintf("%s%s", "http://", path.Join(addr, route))
}

// sendBatch отправляет батч на сервер url
func sendBatch(batch []metrica.Metrica, url string) {
	if len(batch) == 0 {
		return
	}

	log.Debug().Int("batch_len", len(batch)).Msg("Отправляю метрики")
	err := report.Send(batch, url)
	if err != nil {
		log.Error().Msgf("Попытка отправить метрики завершилась с  ошибками: %v", err)
		return
	}
}

// todo Сделать классную такую мермаид диаграмму со всеми каналами, кто куда как собирает
