package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// loadEnv переписывает глобальные параметры настроек адреса и интервалов отправки и опроса
// если такие назначены

func loadEnv() {
	cfg := struct {
		Addr            string `env:"ADDRESS"`
		ReportInterval  int    `env:"REPORT_INTERVAL"`
		PollingInterval int    `env:"POLLING_INTERVAL"`
	}{}

	env.Parse(&cfg)
	if cfg.Addr != "" {
		*addr = cfg.Addr
		fmt.Println("Адрес сервера был переназначен переменной окружения ADDRESS на ", *addr)
	}
	if cfg.ReportInterval != 0 {
		*reportIntervalSeconds = cfg.ReportInterval
		fmt.Println("Интервал отправки данных переназначен переменной окружения REPORT_INTERVAL на ", *reportIntervalSeconds)
	}
	if cfg.PollingInterval != 0 {
		*pollIntervalSeconds = cfg.PollingInterval
		fmt.Println("Интервал опроса метрик переназначен переменной окружения POLLING_INTERVAL на ", *pollIntervalSeconds)
	}
}
