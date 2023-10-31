package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	//"github.com/octago/sflags/gen/gflag" сделать свой репозиторий и залить его всем в ПР
)

type Agent struct {
	Addr            string `env:"ADDRESS" flag:"~a" desc:"(строка) адрес сервера в формате host:port"`
	ReportInterval  uint   `env:"REPORT_INTERVAL" flag:"~r" desc:"(число, секунды) частота отправки данных на сервер"`
	PollingInterval uint   `env:"POLLING_INTERVAL" flag:"~p" desc:"(число, секунды) частота отпроса метрик"`
	RateLimit       uint   `env:"RATE_LIMIT"`
	Key             Secret `env:"KEY" flag:"~p" desc:"(строка) секретный ключ подписи"`
}

// Parse парсит настройки адреса сервера, и частоты опроса и отправки
// из командной строки и переменных окружения. В приоритете переменные окружения
func (cfg *Agent) Parse(defaults Agent) error {
	// todo
	// сделать репозиторий sflags домашним, чтобы он мог устанавливаться от меня хотя бы
	// github.com/octago/sflags, сейчас там ошибка в go.mod
	// тогда можно будет просто сделать gflag.parse(&cfg)

	flag.UintVar(&cfg.PollingInterval, "p", defaults.PollingInterval, "число, частота опроса метрик")
	flag.UintVar(&cfg.ReportInterval, "r", defaults.ReportInterval, "число, частота отправки данных на сервер")
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "строка, адрес сервера в формате host:port")
	flag.Var(&cfg.Key, "k", "строка, секретный ключ подписи")
	flag.UintVar(&cfg.RateLimit, "l", defaults.RateLimit, "число, максимальное количество исходящих запросов")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	// Валидация
	if cfg.RateLimit == 0 {
		return fmt.Errorf("Ошибка при конфигурировании, количество исходящих соединения не может быть меньше 1")
	}

	log.Info().Msgf("Запущено с настройками %+v", cfg)
	return nil
}

func init() {
	// настраиваем дружелюбный цветастый вывод логгера
	log.Logger = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// добавляет смайлик кота в конец справки
	flag.Usage = func() {
		print("server")
		flag.PrintDefaults()
		fmt.Println("^-^")
	}
}
