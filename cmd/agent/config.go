package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServer                = "localhost:8080"
	defaultPollIntervalSeconds   = 2
	defaultReportIntervalSeconds = 10
)

func init() {
	flag.Usage = func() {
		print("server")
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}

type config struct {
	Addr            string `env:"ADDRESS"`
	ReportInterval  int    `env:"REPORT_INTERVAL"`
	PollingInterval int    `env:"POLLING_INTERVAL"`
}

// configure переписывает глобальные параметры настроек адреса и интервалов отправки и опроса
// если такие назначены
func configure() (cfg config) {
	//togo default config можно тоже объявить конфиг структурой, или передать в функцию!

	cfg.PollingInterval = *flag.Int("p", defaultPollIntervalSeconds, "число, частота опроса метрик")
	cfg.ReportInterval = *flag.Int("r", defaultReportIntervalSeconds, "число, частота отправки данных на сервер")
	cfg.Addr = *flag.String("a", defaultServer, "строка, адрес сервера в формате host:port")

	env.Parse(&cfg)

	return
}
