// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
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
	log.Info().Msgf("Сервер запущен строкой %v", strings.Join(os.Args, " "))

	cfg := config.Server{}
	err := cfg.Parse(defaultConfig)
	if err != nil {
		log.Error().Msgf("Ошибка парсинга конфига: %v", err)
		os.Exit(2)
	}
	fmt.Printf("Получен конфиг %+v \n", cfg)

	rootContext, cancel := context.WithCancel(context.Background()) // это пусть будет просто defer storage.Close

	// создаем хранилище
	s, err := cfg.MakeStorage(rootContext)
	if err != nil {
		log.Error().Msgf("Не удалось создать хранилише: %v", err)
		return
	}

	// создаем роутер
	router := router.MeowRouter(s, string(cfg.Key.ValueFunc()()))

	// Запускаем сервер с поддержкой нежного завершения,
	// занимаем текущий поток до вызова сигналов выключения
	graceful.Serve(cfg.Addr, router)

	cancel()

	// Даем ему время
	time.Sleep(time.Second)

	log.Info().Msg("^.^ Сервер завершен нежно")
	// Wait for server context to be stopped

}
