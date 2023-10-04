package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	//"github.com/octago/sflags/gen/gflag" сделать свой репозиторий и залить его всем в ПР
)

type config struct {
	Addr            string `env:"ADDRESS" flag:"~a" desc:"(строка) адрес сервера в формате host:port"`
	ReportInterval  uint   `env:"REPORT_INTERVAL" flag:"~r" desc:"(число, секунды) частота отправки данных на сервер"`
	PollingInterval uint   `env:"POLLING_INTERVAL" flag:"~p" desc:"(число, секунды) частота отпроса метрик"`
}

var defaultConfig = config{
	Addr:            "localhost:8080",
	ReportInterval:  10,
	PollingInterval: 2,
}

// mustConfigure парсит настройки адреса сервера, и частоты опроса и отправки
// из командной строки и переменных окружения. В приоритете переменные окружения
func mustConfigure(defaults config) (cfg config) {
	// todo
	// сделать репозиторий sflags домашним, чтобы он мог устанавливаться от меня хотя бы
	// github.com/octago/sflags, сейчас там ошибка в go.mod
	// тогда можно будет просто сделать gflag.parse(&cfg)

	flag.UintVar(&cfg.PollingInterval, "p", defaults.PollingInterval, "число, частота опроса метрик")
	flag.UintVar(&cfg.ReportInterval, "r", defaults.ReportInterval, "число, частота отправки данных на сервер")
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "строка, адрес сервера в формате host:port")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	ololog.Info().Msgf("Запущено с настройками %+v", cfg)

	return
}

func init() {
	// добавляет смайлик кота в конец справки
	flag.Usage = func() {
		print("server")
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
