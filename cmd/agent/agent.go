package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/report"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/report/compress"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
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
	// подготовим выключение
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Запускаем работу
	// после получения сигнала агент отправит последние
	// метрики и завершится
	collector.FetchAndReport(ctx, config, updateRoute)

	// Все завершилось, выведем последнюю статистику
	log.Info().Int("goroutines active", runtime.NumGoroutine()).Msgf("Работа завершена нежно")
}
