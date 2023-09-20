package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	//"github.com/octago/sflags/gen/gflag" сделать свой репозиторий и залить его всем в ПР
)

type config struct {
	Addr            string `env:"ADDRESS" flag:"~a" desc:"(строка) адрес сервера в формате host:port"`
	ReportInterval  int    `env:"REPORT_INTERVAL" flag:"~r" desc:"(число, секунды) частота отправки данных на сервер"`
	PollingInterval int    `env:"POLLING_INTERVAL" flag:"~p" desc:"(число, секунды) частота отпроса метрик"`
}

// configure переписывает глобальные параметры настроек адреса и интервалов отправки и опроса
// если такие назначены
func configure(defaults config) (cfg config) {
	// todo
	// сделать репозиторий sflags домашним, чтобы он мог устанавливаться от меня хотя бы
	// github.com/octago/sflags, сейчас там ошибка в go.mod
	// тогда можно будет просто сделать gflag.parse(&cfg)

	flag.IntVar(&cfg.PollingInterval, "p", defaults.PollingInterval, "число, частота опроса метрик")
	flag.IntVar(&cfg.ReportInterval, "r", defaults.ReportInterval, "число, частота отправки данных на сервер")
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "строка, адрес сервера в формате host:port")

	flag.Parse()

	env.Parse(&cfg)

	return
}

func init() {
	flag.Usage = func() {
		print("server")
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
