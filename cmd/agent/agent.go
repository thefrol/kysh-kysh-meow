package main

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/compress"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/report"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
	"github.com/thefrol/kysh-kysh-meow/lib/graceful"
)

const (
	updateRoute            = "/updates"
	GracefulShutdownPeriod = 30 * time.Second
	CompressMinLength      = 100 // порог байтов после которого начинаем сжимать ответ
)

var defaultConfig = config.Agent{
	Addr:            "localhost:8080",
	ReportInterval:  10,
	PollingInterval: 2,
	RateLimit:       2,
}

func main() {
	// Парсим командную строку и переменные окружения
	config := config.Agent{}
	if err := config.Parse(defaultConfig); err != nil {
		log.Error().Msgf("Ошибка парсинга конфига: %v", err)
		os.Exit(2)
	}

	// Настроим отправку

	// todo Вообще все как-то почистить это все в агенте, например
	// руководствуясб https://www.digitalocean.com/community/tutorials/understanding-init-in-go-ru

	// MENTOR я вообще не оч хорошо понимаю,
	// нормально ли что это вообще какой-то левый пакет настривается тут?
	// возможно все, вплоть до UpdateRoute должно уехать в конфиг, даже
	// если такие переменные и не связаны с командной строкой. Зато красиво
	// в одной месте можно оформить установку всех основных параметров
	report.SetSigningKey(string(config.Key.ValueFunc()()))
	report.CompressLevel = compress.BestCompression
	report.CompressMinLength = CompressMinLength

	//Создадим контекст, который будет завершен по сигналу ОС
	ctx := graceful.WithSignal(context.Background())
	graceful.SetForcedShutdown(ctx, GracefulShutdownPeriod)

	// Запускаем работу
	collector.FetchAndReport(ctx, config, updateRoute)

	// Все завершилось, выведем последнюю статистику
	log.Info().Int("goroutines active", runtime.NumGoroutine()).Msgf("Работа завершена нежно")

	// todo на данный момент в конце работают три горутины:
	// + текущая(main())
	// + горутина graceful.ForcedShutdown() - при желании можно остановить
	// + горутина signal.Notify() при желании можно тоже остановить, наверно
	//
	// То есть, 3 оставшиеся горутины значит, что все эти бесконечные остальные
	// горутины нормально у меня закрываются
}

// todo
//
// было бы прикольно ещё придумать способ, как бы сделать так, чтобы
// report.Send не занимал семафоры, пока идёт ожидание отправки,
// например, туда можно было бы передавать этот семафор как раз.
// Но тогда семаформы и всякие примитивы понадобится выделить
// в отдельный пакет
