// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"

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

	// создаем хранилище
	s, cancelStorage := cfg.MakeStorage()
	// Запускаем сервер с поддержкой нежного завершения,
	// занимаем текущий поток до вызова сигнатов выключения
	Run(cfg, s) // будет app.Run()

	// Завершаем последние дела
	// попытаемся сохраниться в файл
	cancelStorage()

	// Даем ему время
	time.Sleep(time.Second)

	log.Info().Msg("^.^ Сервер завершен нежно")
	// Wait for server context to be stopped

}

// Run запускает сервер с поддержкой нежного завершения. Сервер можно будет выключить через
// SIGINT, SIGTERM, SIGQUIT
func Run(cfg config.Server, s api.Operator) {
	// Запускаем сервер с поддержкой нежного выключения
	// вдохноввлено примерами роутера chi
	server := http.Server{Addr: cfg.Addr, Handler: router.MeowRouter(s, string(cfg.Key.ValueFunc()()))}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Debug().Msg("signal received")
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
				return
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Msg(err.Error())
			panic(err)
		}
		serverStopCtx()
		log.Info().Msg("^-^ рутина остановки сервера завершилась")
	}()
	log.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", cfg.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Msgf("^0^ не могу запустить сервер: %v \n", err)
		//todo
		//
		// если не биндится, то хотя бы выходить с ошибкой,
		// в данный момент сервер не закроета сам
		//
		// можно дать несколько попыток забиндиться
	}

	<-serverCtx.Done()

	// МНе кажется в отдельную функцию надо выделить именно все, что относится к нежному завершению, + надо перевести комменты по коду на русский
}
