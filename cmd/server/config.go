package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	Addr                 string `env:"ADDRESS"`
	StoreIntervalSeconds uint   `env:"STORE_INTERVAL"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	Restore              bool   `env:"RESTORE"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
	Key                  string `env:"KEY"`
}

var defaultConfig = config{
	Addr:                 ":8080",
	StoreIntervalSeconds: 300,
	FileStoragePath:      "/tmp/metrics-db.json",
	Restore:              true, // в текущей конфигурации это значение командной строкой никак не поменять, нельзя указать -r 0, флан такое не принимает todo
}

// mustConfigure парсит командную строку и переменные окружения, чтобы выдать структуру с конфигурацией сервера.
// В приоритете переменные окружения. Принимает на вход структуру defaults со значениями по умолчанию.
//
// Приоритет такой:
//   - Если другого не указано, будет использоваться defaults
//   - То, что указано в командной строке переписывает то, что указано в defaults
//   - То, что указано в переменной окружения, переписывает то, что было указано ранее
func mustConfigure(defaults config) (cfg config) {
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "[адрес:порт] устанавливает адрес сервера ")
	flag.UintVar(&cfg.StoreIntervalSeconds, "i", defaults.StoreIntervalSeconds, "[время, сек] интервал сохранения показаний. При 0 запись делается почти синхронно")
	flag.StringVar(&cfg.FileStoragePath, "f", defaults.FileStoragePath, "[строка] путь к файлу, откуда будут читаться при запуске и куда будут сохраняться метрики полученные сервером, если файл пустой, то сохранение будет отменено")
	flag.BoolVar(&cfg.Restore, "r", defaultConfig.Restore, "[флаг] если установлен, загружает из файла ранее записанные метрики")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaults.DatabaseDSN, "[строка] подключения к базе данных")
	flag.StringVar(&cfg.Key, "k", defaults.Key, "строка, секретный ключ подписи")

	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	// Тут обрабатываем особый случай. Если переменная окружения установлена, но в пустое значение
	// то мы перезаписываем установленный командной строкой флаг на пуское значение, хотя штатно
	// этого не было бы сделано
	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = v
	}

	return
}

func init() {
	// настраиваем дружелюбный цветастый вывод логгера
	log.Logger = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// добавляет смайлик кота в конец справки flags
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
