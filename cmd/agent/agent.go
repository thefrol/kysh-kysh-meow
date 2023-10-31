package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/report"
	"github.com/thefrol/kysh-kysh-meow/internal/compress"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
)

const updateRoute = "/updates"

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

	// MENTOR я вообще не оч хорошо понимаю,
	// нормально ли что это вообще какой-то левый пакет настривается тут?
	// возможно все, вплоть до UpdateRoute должно уехать в конфиг, даже
	// если такие переменные и не связаны с командной строкой. Зато красиво
	// в одной месте можно оформить установку всех основных параметров
	report.SetSigningKey(string(config.Key.ValueFunc()()))
	report.CompressLevel = compress.BestCompression
	report.CompressMinLength = 100

	// Запускаем работу
	collector.FetchAndReport(config, updateRoute)
}
