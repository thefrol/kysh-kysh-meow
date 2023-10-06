// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"
)

func main() {
	cfg := mustConfigure(defaultConfig)

	// создаем хранилище
	s, err := fileStorage(cfg)
	if err != nil {
		ololog.Error().Msgf("Не удалось сконфигурировать сервер, по причине: %v", err)
		return
	}

	// Запускаем сервер с поддержкой нежного выключения
	// вдохноввлено примерами роутера chi
	server := http.Server{Addr: cfg.Addr, Handler: router.MeowRouter(s)}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		ololog.Debug().Msg("signal received")
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
	ololog.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", cfg.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		ololog.Error().Msgf("^0^ не могу запустить сервер: %v \n", err)
	}

	// Завершаем последние дела
	// попытаемся сохраниться в файл

	// saver позволяет получить доступ к функции сохранения в хранилище,
	// потому что в Storager мы его не имеем, и не очень-то хотим
	if cfg.FileStoragePath != "" { //todo можно сделать как функцию в структуре config
		type saver interface {
			ToFile(string) error
		}

		if v, ok := s.(saver); ok {
			v.ToFile(cfg.FileStoragePath)
			ololog.Info().Msg("Сохранено в файл")
		} else {
			ololog.Error().Msg("Не могу преобразовать хранилище в нужный интерфейс для сохранения данных на выходе")
		}
	}

	ololog.Info().Msg("^.^ Сервер завершен нежно")
	// Wait for server context to be stopped
	<-serverCtx.Done()

}
