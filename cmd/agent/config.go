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

var defaultConfig = config{
	Addr:            defaultServer,
	ReportInterval:  defaultReportIntervalSeconds,
	PollingInterval: defaultPollIntervalSeconds,
}

// configure переписывает глобальные параметры настроек адреса и интервалов отправки и опроса
// если такие назначены
func configure(defaults config) (cfg config) {
	//togo default config можно тоже объявить конфиг структурой, или передать в функцию!
	flag.IntVar(&cfg.PollingInterval, "p", defaults.PollingInterval, "число, частота опроса метрик")
	flag.IntVar(&cfg.ReportInterval, "r", defaults.ReportInterval, "число, частота отправки данных на сервер")
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "строка, адрес сервера в формате host:port")

	flag.Parse()

	fmt.Printf("in configure()=%+v\n", cfg)

	env.Parse(&cfg)

	return
}
