package config

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router/httpio"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
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

// ConfigureStorage подготавливает хранилище к работе в соответствии с текущими настройками,
// при необходимости загружает из файла значения метрик, запускает сохранение в файл, и
// возвращает интерфейс хранилища и функцию, подготавливающая ханилище к остановке
//
// На входе получает экземпляр хранилища m, и далее оборачивает его другим классов,
// наиболее соответсвующим задаче, исходя из cfg
func (cfg Server) MakeStorage(ctx context.Context) (httpio.Operator, error) {
	// todo эта функция должна переехать в пакет storage

	// 0. Если указана база данных, создаем хранилище с базой данных
	// 1. Если путь не задан, то возвращаем хранилище в оперативке, без приблуд
	// 2. Иначе оборачиваем файловым хранилищем, но не возвращаем пока
	// 3. Если Restore=true, то читаем из файла. Если файла не существует, то игнорируем проблему
	// 4. Оборачиваем файловое хранилище сихнронным или интервальным

	// TODO
	//
	// Думаю эта функция должна получать на вход контекст, а не возвращать CancelFunc
	//
	// И нужно убрать отсюда все паники

	// Если база данных
	if cfg.DatabaseDSN.Get() != "" {
		db, err := sql.Open("pgx", cfg.DatabaseDSN.Get())
		if err != nil {
			return nil, fmt.Errorf("не могу создать соединение с БД: %w", err)
		}

		dbs, err := storage.NewDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания хранилища в базе данных: %v", err)
		}

		if err := dbs.Check(context.TODO()); err != nil {
			log.Warn().Msgf("Нет соединения с БД - %v", err)
		}

		// соединения с базой данных закроются по завершении работы
		// не обязательно явно это делать

		log.Info().Msg("Создано хранилише в Базе данных")
		return dbs, nil
	}

	// Если не база данных, то начинаем с начала - создаем хранилище в памяти, и оборачиваем его всякими штучками если надо
	m := storage.New()

	// Если путь до хранилища не пустой, то нам нужно инициаизировать обертки над хранилищем
	if cfg.FileStoragePath == "" {
		log.Info().Msg("Установлено хранилище в памяти. Сохранение на диск отключено")
		return storage.AsOperator(m), nil
	}

	// Оборачиваем файловым хранилищем, в случае, есл и
	fs := storage.NewFileStorage(&m, cfg.FileStoragePath)
	if cfg.Restore {
		err := fs.Restore()
		// в случае, если файла не существует, игнорируем эту проблему
		if err != nil && err != storage.ErrorRestoreFileNotExist {
			panic(err)
		}
		log.Info().Msgf("Значения метрик загружены из %v", fs.FileName)
	}

	if cfg.StoreIntervalSeconds == 0 {
		// инициализируем сихнронную запись,
		// при этом сохраняться в конце нам не понадобится
		log.Info().Msgf("Установлено синхронное сохранение в %v в при записи", fs.FileName)
		return storage.AsOperator(storage.NewSyncDump(&fs)), nil
	}

	// Запускаем интервальную запись
	s := storage.NewIntervalDump(&fs, time.Duration(cfg.StoreIntervalSeconds)*time.Second)
	go s.StartDumping(ctx)

	log.Info().Msgf("Установлено сохранение с интервалом %v в %v в при записи", s.Interval, s.FileName)

	return storage.AsOperator(s), nil
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

// TODO может быть storage должно иметь что-то типа Close()?
