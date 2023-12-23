package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Addr                 string           `env:"ADDRESS"`
	StoreIntervalSeconds uint             `env:"STORE_INTERVAL"`
	FileStoragePath      string           `env:"FILE_STORAGE_PATH"`
	Restore              bool             `env:"RESTORE"`
	DatabaseDSN          ConnectionString `env:"DATABASE_DSN"`
	Key                  Secret           `env:"KEY"`
}

// Parse парсит командную строку и переменные окружения, чтобы выдать структуру с конфигурацией сервера.
// В приоритете переменные окружения. Принимает на вход структуру defaults со значениями по умолчанию.
//
// Приоритет такой:
//   - Если другого не указано, будет использоваться defaults
//   - То, что указано в командной строке переписывает то, что указано в defaults
//   - То, что указано в переменной окружения, переписывает то, что было указано ранее
func (cfg *Server) Parse(defaults Server) error {
	// устанавливаем дополнительные значения по умолчанию
	cfg.DatabaseDSN.s = defaults.DatabaseDSN.s

	// парсим командную строку
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "[адрес:порт] устанавливает адрес сервера ")
	flag.UintVar(&cfg.StoreIntervalSeconds, "i", defaults.StoreIntervalSeconds, "[время, сек] интервал сохранения показаний. При 0 запись делается почти синхронно")
	flag.StringVar(&cfg.FileStoragePath, "f", defaults.FileStoragePath, "[строка] путь к файлу, откуда будут читаться при запуске и куда будут сохраняться метрики полученные сервером, если файл пустой, то сохранение будет отменено")
	flag.BoolVar(&cfg.Restore, "r", defaults.Restore, "[флаг] если установлен, загружает из файла ранее записанные метрики")
	flag.Var(&cfg.DatabaseDSN, "d", "[строка] подключения к базе данных")
	flag.Var(&cfg.Key, "k", "строка, секретный ключ подписи")

	flag.Parse()
	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	// Тут обрабатываем особый случай. Если переменная окружения установлена, но в пустое значение
	// то мы перезаписываем установленный командной строкой флаг на пуское значение, хотя штатно
	// этого не было бы сделано
	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = v
	}

	log.Info().Msgf("Запущено с настройками %+v", cfg)
	return nil
}
