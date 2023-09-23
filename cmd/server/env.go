package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// loadEnv переписывает глобальный параметр настроек адреса сервера
func loadEnv() {
	cfg := struct {
		Addr string `env:"ADDRESS"`
	}{}

	env.Parse(&cfg)
	if cfg.Addr != "" {
		*addr = cfg.Addr
		fmt.Println("Адрес сервера был переписан переменной окружения ADDRESS на ", *addr)
	}
}
