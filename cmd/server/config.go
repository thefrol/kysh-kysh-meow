package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

func init() {

	flag.Usage = func() {
		print("server")
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}

type config struct {
	Addr string `env:"ADDRESS"`
}

// configure парсит командную строку и переменные окружения, чтобы выдать структуру с конфигурацией сервера
func configure(defaultServer string) (cfg config) {
	flag.StringVar(&cfg.Addr, "a", defaultServer, "[адрес:порт] устанавливает адрес сервера ")
	env.Parse(&cfg)
	return
}
