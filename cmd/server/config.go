package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

type config struct {
	Addr                 string `env:"ADDRESS"`
	StoreIntervalSeconds uint   `env:"STORE_INTERVAL"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	Restore              bool   `env:"RESTORE"`
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

	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	// todo
	//
	// вообще две эти функции сверху требуют проверку ошибок, и это в тестах тоже стоило бы отразить

	// Тут обрабатываем особый случай. Если переменная окружения установлена, но в пустое значение
	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = v
	}

	ololog.Info().Msgf("Запущено с настройками %+v", cfg)

	return
}

func init() {
	// добавляет смайлик кота в конец справки flags
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
