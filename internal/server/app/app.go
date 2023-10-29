package app

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type App struct {
	DB sql.Conn
}

// Run запускает сервер с поддержкой нежного завершения. Сервер можно будет выключить через
// SIGINT, SIGTERM, SIGQUIT
func (a *App) Run(addr string, router http.Handler) {
	// Запускаем сервер с поддержкой нежного выключения
	// вдохноввлено примерами роутера chi
	server := http.Server{Addr: addr, Handler: router}

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
	log.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", addr)

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
	log.Error().Msg("Run() остановлен")

	// МНе кажется в отдельную функцию надо выделить именно все, что относится к нежному завершению, + надо перевести комменты по коду на русский
}
