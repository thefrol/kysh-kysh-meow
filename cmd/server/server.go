// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge
package main

import (
	"context"
	"database/sql"
	"os"

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

	// storageContext это контекст бд, он нужен чтобы
	// остановить горутины, который пишут в файл, а
	// так же закрыть соединение с бд при остановку сервера
	storageContext, stopStorage := context.WithCancel(
		context.Background())

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

		counters = &s
		gauges = &s
		labels = &s

	} else {

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
	log.Info().Msg("Сервер остановлен")

	// Останавливаем хранилище, интервальную запись в файл и все остальное
	// или соединения с БД
	log.Info().Msg("Останавливаем хранилище")
	stopStorage()

	log.Info().Msg("^.^ Сервер завершен нежно")
	// Wait for server context to be stopped
}
