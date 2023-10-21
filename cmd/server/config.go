package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Addr string `env:"ADDRESS"`
}

// configure парсит командную строку и переменные окружения, чтобы выдать структуру с конфигурацией сервера.
// В приоритете переменные окружения,
func configure(defaults config) (cfg config) {
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "[адрес:порт] устанавливает адрес сервера ")

	flag.Parse()
	env.Parse(&cfg)
	return
}

func init() {
	// добавляет смайлик кота в конец справки
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
