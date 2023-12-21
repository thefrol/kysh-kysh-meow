// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge
package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/dbping"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/metricas"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/scan"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"
	"github.com/thefrol/kysh-kysh-meow/internal/server/storage"
	"github.com/thefrol/kysh-kysh-meow/internal/server/storagev2/mem"
	"github.com/thefrol/kysh-kysh-meow/lib/graceful"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var defaultConfig = config.Server{
	Addr:                 ":8080",
	StoreIntervalSeconds: 300,
	FileStoragePath:      "/tmp/metrics-db.json",
	Restore:              true, // в текущей конфигурации это значение командной строкой никак не поменять, нельзя указать -r 0, флан такое не принимает todo
}

func main() {
	// Парсим командную строку и переменные окружения
	cfg := config.Server{}
	if err := cfg.Parse(defaultConfig); err != nil {
		log.Error().Msgf("Ошибка парсинга конфига: %v", err)
		os.Exit(2)
	}

	// создаем хранилище
	// сейчас мы создадим его так, что
	// 1. Если dsn не указал, то используем
	// storagev2
	// 2. Иначе, исползуем старый вариант

	var (
		counters manager.CounterRepository
		gauges   manager.GaugeRepository
		labels   scan.Labler
	)

	if cfg.DatabaseDSN.Get() == "" {
		log.Info().Msg("используется storagev2")

		s := mem.MemStore{
			Counters: make(mem.IntMap, 50),
			Gauges:   make(mem.FloatMap, 50),

			Log: log.With().Str("storage", "memStoreV2").Logger(),
		}

		// если указан флаг restore, то читаем из нашего
		// файла
		if cfg.Restore {
			err := s.RestoreFrom(cfg.FileStoragePath)
			if err != nil {
				log.Error().
					Err(err).
					Msg("не удалось загрузить storage из файла")

				os.Exit(1)

				//todo
				//
				// resore заменит уже созданные мапы, на те, что куче
				// надо как-то копировать их, а не просто заменять
			}

		}

		// теперь разберемся с сохранением,
		// если интевал нулевой - делае					fmt.Print("123")м сохранение
		// синхронным, силами структуры
		if cfg.StoreIntervalSeconds == 0 {
			s.FilePath = cfg.FileStoragePath
		} else {

			// если мы используем интервальное сохранение
			// то у нас для этого есть специальный класс
			i := mem.IntervalicSaver{
				Store:    &s,
				File:     cfg.FileStoragePath,
				Interval: time.Duration(cfg.StoreIntervalSeconds) * time.Second,
			}

			// Запускаем интервальное сохранение
			// и на выходе из мейна поставим ожидание
			err := i.Run()
			if err != nil {
				log.Error().
					Err(err).
					Msg("не запущено интервальное сохранение")
				os.Exit(1)
			}
			defer i.Stop()

		}
		counters = &s
		gauges = &s
		labels = &s
	} else {

		// storageContext это уже устаревшая тема, нам не нужен
		// какой-то крейзи контекст БД
		storageContext, stopStorage := context.WithCancel(
			context.Background())
		defer stopStorage()

		s, err := cfg.MakeStorage(storageContext)
		if err != nil {
			log.Error().Msgf("Не удалось создать хранилише: %v", err)
			return
		}

		// готовим репозитории
		counters = &storage.CounterAdapter{
			Op: s,
		}

		gauges = &storage.GaugeAdapter{
			Op: s,
		}

		labels = &storage.LabelsAdapter{
			Op: s,
		}
	}
	// готовим прикладной уровень
	scanner := scan.Labels{
		Labels: labels,
	}

	reg := manager.Registry{
		Counters: counters,
		Gauges:   gauges,
	}

	man := metricas.Manager{
		Registry: reg,
	}

	// создаем пингер
	//
	// у него будет свое соединение!
	db, err := sql.Open("pgx", cfg.DatabaseDSN.Get())
	if err != nil {
		log.Error().Msgf("не могу создать соединение с БД: %v", err)
		os.Exit(1)
	}

	// вообще пингер загадочная штука, поэтому у него свой юзкейс
	pinger := dbping.Pinger{
		Connection: db,
	}

	// создаем роутер
	r := router.API{
		Manager:   man,
		Registry:  reg,
		Dashboard: scanner,
		Pinger:    pinger,

		Key: string(cfg.Key.ValueFunc()()), // todo это лол
	}

	router := r.MeowRouter()

	// Запускаем сервер с поддержкой нежного завершения,
	// занимаем текущий поток до вызова сигналов выключения
	err = graceful.ListenAndServe(context.Background(), cfg.Addr, router)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка запуска сервера")
	}
	log.Info().Msg("Апи остановлен")

	// конец. парам па-па пам
	log.Info().Msg("^.^ Сервер завершен нежно, остались деферы")

}
