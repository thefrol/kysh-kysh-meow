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
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

type Server struct {
	Addr                 string `env:"ADDRESS"`
	StoreIntervalSeconds uint   `env:"STORE_INTERVAL"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	Restore              bool   `env:"RESTORE"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
	Key                  Secret `env:"KEY"`
}

// mustConfigure парсит командную строку и переменные окружения, чтобы выдать структуру с конфигурацией сервера.
// В приоритете переменные окружения. Принимает на вход структуру defaults со значениями по умолчанию.
//
// Приоритет такой:
//   - Если другого не указано, будет использоваться defaults
//   - То, что указано в командной строке переписывает то, что указано в defaults
//   - То, что указано в переменной окружения, переписывает то, что было указано ранее
func (cfg *Server) Parse(defaults Server) error {
	flag.StringVar(&cfg.Addr, "a", defaults.Addr, "[адрес:порт] устанавливает адрес сервера ")
	flag.UintVar(&cfg.StoreIntervalSeconds, "i", defaults.StoreIntervalSeconds, "[время, сек] интервал сохранения показаний. При 0 запись делается почти синхронно")
	flag.StringVar(&cfg.FileStoragePath, "f", defaults.FileStoragePath, "[строка] путь к файлу, откуда будут читаться при запуске и куда будут сохраняться метрики полученные сервером, если файл пустой, то сохранение будет отменено")
	flag.BoolVar(&cfg.Restore, "r", defaults.Restore, "[флаг] если установлен, загружает из файла ранее записанные метрики")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaults.DatabaseDSN, "[строка] подключения к базе данных")
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
	return nil
}

// ConfigureStorage подготавливает хранилище к работе в соответствии с текущими настройками,
// при необходимости загружает из файла значения метрик, запускает сохранение в файл, и
// возвращает интерфейс хранилища и функцию, подготавливающая ханилище к остановке
//
// На входе получает экземпляр хранилища m, и далее оборачивает его другим классов,
// наиболее соответсвующим задаче, исходя из cfg
func (cfg Server) MakeStorage() (api.Operator, context.CancelFunc, error) {
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
	if cfg.DatabaseDSN != "" {
		db, err := sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			return nil, nil, fmt.Errorf("не могу создать соединение с БД: %w", err)
		}

		dbs, err := storage.NewDatabase(db)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка создания хранилища в базе данных: %v", err)
		}

		if err := dbs.Check(context.TODO()); err != nil {
			log.Warn().Msgf("Нет соединения с БД - %v", err)
		}

		log.Info().Msg("Создано хранилише в Базе данных")
		return dbs, func() {
			err := db.Close()
			if err != nil {
				log.Error().Msgf("Не могу закрыть базу данных: %v", err)
			}

			// todo
			//
			// Конечно, я хочу делать это defer или как-то так, можно у нас будет некий app.Close()
		}, nil
	}

	// Если не база данных, то начинаем с начала - создаем хранилище в памяти, и оборачиваем его всякими штучками если надо
	m := storage.New()

	// Если путь до хранилища не пустой, то нам нужно инициаизировать обертки над хранилищем
	if cfg.FileStoragePath == "" {
		log.Info().Msg("Установлено хранилище в памяти. Сохранение на диск отключено")
		return storage.AsOperator(m), func() {
			log.Info().Msg("Хранилище сихронной записи получило сигнал о завершении, но файловая запись в текущей конфигурации сервера не используется. Ничего не записано")
		}, nil
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
		return storage.AsOperator(storage.NewSyncDump(&fs)), func() {
			log.Info().Msg("Хранилище сихронной записи получило сигнал о завершении, но дополнительно сохранение не нужно")
		}, nil
	}

	// Запускаем интервальную запись и создаем токен отмены, при необходимости сюда можно будет добавить и группу ожидания
	s := storage.NewIntervalDump(&fs, time.Duration(cfg.StoreIntervalSeconds)*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go s.StartDumping(ctx)

	log.Info().Msgf("Установлено сохранение с интервалом %v в %v в при записи", s.Interval, s.FileName)

	return storage.AsOperator(s), func() {
		log.Info().Msg("Хранилище интервальной записи получило сигнал о завершении, но дополнительно сохранение не нужно")
		cancel()
	}, nil

	// TODO
	//
	// Пока что судя по всему эта функция эвакуирует моё хранилище из стека, мне кажется, но не мапы
	// Как вариант сразу создавать FileStorage, и просто не оборачивать его если надо
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
