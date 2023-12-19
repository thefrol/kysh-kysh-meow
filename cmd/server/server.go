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
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"
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
	s, err := cfg.MakeStorage(storageContext)
	if err != nil {
		log.Error().Msgf("Не удалось создать хранилише: %v", err)
		return
	}

	// создаем пингер
	//
	// у него будет свое соединение!
	db, err := sql.Open("pgx", cfg.DatabaseDSN.Get())
	if err != nil {
		log.Error().Msgf("не могу создать соединение с БД: %v", err)
		os.Exit(1)
	}

	pinger := dbping.Pinger{
		Connection: db,
	}

	// создаем роутер
	router := router.MeowRouter(s, pinger, string(cfg.Key.ValueFunc()()))

	// Запускаем сервер с поддержкой нежного завершения,
	// занимаем текущий поток до вызова сигналов выключения
	graceful.Serve(cfg.Addr, router)

	// Останавливаем хранилище, интервальную запись в файл и все остальное
	// или соединения с БД
	stopStorage()

	// Даем ему время
	time.Sleep(time.Second)

	log.Info().Msg("^.^ Сервер завершен нежно")
	// Wait for server context to be stopped

}
