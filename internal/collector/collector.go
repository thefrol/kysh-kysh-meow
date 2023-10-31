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

func FetchAndReport(ctx context.Context, config config.Agent, updateRoute string) {

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
	inMs := generator(ctx, fetch.MemStats, interval)
	inPc := generator(ctx, fetch.PollCount, interval)
	inRv := generator(ctx, fetch.RandomValue, interval)
	inPs := generator(ctx, fetch.GoPS, interval)

	// объединить каналы в один
	inMix := FanIn(ctx, inMs, inPc, inPs, inRv)

	// собирать данные с перерывами
	reportInterval := time.Second * time.Duration(config.ReportInterval)
	inCh := WithTimeouts(inMix, reportInterval)

	// отправим данные
	workerCount := 3
	url := Endpoint(config.Addr, updateRoute)
	sema := NewSemaphore(int(config.RateLimit))
	for i := 0; i < workerCount; i++ {
		worker(inCh, url, sema)
	}
	// надо конечно подумать над такими вещами, что если метрики не отправились? бросить их обратно в начало или в какую-то
	// Дополнительную очередь?
	// типа у нас ещё очеред в начало одна, мертвых сообщений они опять на входе

	// что делать если батч ещё не собран до конца, чтобы отправлять
	// тут сильно поможет решение с созданием каналов
	// в таком случа мы можем собирать какой угодно большой батч

	wg.Wait()

	// todo
	//
	// я наконец конял как должен работать поллкаунт. Он не привязан
	// к остальным метрикам. Типа если бы у меня был каунтер Overflow
	// количество сообщений о переполнении, то я бы хотел чтобы просто
	// копилось количество раз сколько это случилось. Мне важно, чтобы я
	// я видел сколько это случилсоь за то время что я наблюдаю.
	// Вне зависимости от сториджа! Это никак не связано с остальными метриками.
	// Пусть они вообще хоть отдельной гоурутиной приходят.
	// С поллкаунт также. Надо вообще убрать такую штуку как дроп каунтер.
	// и инкремент тоже! (!!! bug)
	//
	//
	// А ещё тут бы я сделал структуру Package. в которую можно добавлять метрики
	// или функцию Squash(pack []Metrica), которая просуммирует все каунтеры
	// и заменит все гауж на более свежие. единственное, что надо бы знать какой из
	// гаужей самый свежий ахах
}

// Endpoint формирует точку, куда агент будет посылать все запросы на основе своей текущей конфигурации
func Endpoint(addr, route string) string {
	return fmt.Sprintf("%s%s", "http://", path.Join(addr, route))
}

// sendBatch отправляет батч на сервер url
func sendBatch(batch []metrica.Metrica, url string) {
	// отправляем

	// todo
	//
	// проверить что батч не пустой

	log.Debug().Int("batch_len", len(batch)).Msg("Отправляю метрики")
	err := report.Send(batch, url)
	if err != nil {
		log.Error().Msgf("Попытка отправить метрики завершилась с  ошибками: %v", err)
		return
	}

	// todo по хорошему именно вот эту штуку бы в горутине обрабатывать, пусть именно она в пуле висит
	// а мы раз в десять секунд собираем в батчи и передаем дальше.
	// типа добавим ещё один элемент конвеера

}
