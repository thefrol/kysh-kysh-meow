package main

import (
	"flag"
	"fmt"
)

const (
	defaultServer                = "localhost:8080"
	defaultPollIntervalSeconds   = 2
	defaultReportIntervalSeconds = 10
)

var (
	pollIntervalSeconds   = flag.Int("p", defaultPollIntervalSeconds, "число, частота опроса метрик")
	reportIntervalSeconds = flag.Int("r", defaultReportIntervalSeconds, "число, частота отправки данных на сервер")
	addr                  = flag.String("a", defaultServer, "строка, адрес сервера в формате host:port")
)

func init() {
	flag.Usage = func() {
		print("server")
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
