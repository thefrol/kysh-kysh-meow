package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// Serve запускает сервер с поддержкой нежного завершения. Сервер можно будет выключить через
// SIGINT, SIGTERM, SIGQUIT
func Serve(addr string, router http.Handler) {
	// Запускаем сервер с поддержкой нежного выключения
	// вдохноввлено примерами роутера chi
	server := http.Server{Addr: addr, Handler: router}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		log.Debug().Msg("server wants to shut down")
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
		log.Info().Msg("^-^ рутина остановки сервера завершилась")
	}()

	log.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Msgf("^0^ не могу запустить сервер: %v \n", err)
	}
	log.Error().Msg("Run() остановлен")

	// не понимаю в таком случае вырубится ли тот потом что должен был сигналы ос
}
