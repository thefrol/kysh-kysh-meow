package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Addr                 string `env:"ADDRESS"`
	StoreIntervalSeconds int    `env:"STORE_INTERVAL"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	Restore              bool   `env:"RESTORE"`
}

// configure парсит командную строку и переменные окружения, чтобы выдать структуру с конфигурацией сервера.
// В приоритете переменные окружения,
func configure(defaults config) (cfg config) {
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "[адрес:порт] устанавливает адрес сервера ")
	flag.IntVar(&cfg.StoreIntervalSeconds, "i", defaults.StoreIntervalSeconds, "[время, сек] интервал сохранения показаний. При 0 запись делается почти синхронно")
	flag.StringVar(&cfg.FileStoragePath, "f", defaults.FileStoragePath, "[строка] путь к файлу, откуда будут читаться при запуске и куда будут сохраняться метрики полученные сервером, если файл пустой, то сохранение будет отменено")
	flag.BoolVar(&cfg.Restore, "r", defaultConfig.Restore, "[флаг] если установлен, загружает из файла ранее записанные метрики")

	flag.Parse()
	env.Parse(&cfg)

	// todo
	//
	// вообще две эти функции сверху требуют проверку ошибок, и это в тестах тоже стоило бы отразить

	// Тут обрабатываем особый случай. Если переменная окружения установлена, но в пустое значение
	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = v
	}
	return
}

func init() {
	// добавляет смайлик кота в конец справки flags
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
