// Пакет для аккуратного выключения сервера
package graceful

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const timeout = 30 * time.Second

// ListenAndServe копирует функциональность http.ListenAndServe
// но только с завершением по сигналам операционной системы, таким как:
// syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT
//
// Позволяет задать базовый контекст, который работает так: если ctx
// выключится, это спровоцирует аккуратную остановку сервера.
func ListenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	// сделаем сервер
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// подготовим выключение
	ctx, stop := signal.NotifyContext(ctx,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// запускаем горутину, которая будет ждать и получения сигнала на остановку
	// из контекста ctx. При этом нам очень важно дождаться завершения
	// server.Shutdown() - после этого сервер точно остановился
	eg := errgroup.Group{}
	eg.Go(func() error {
		// ждем, что придет сигнал на остановку
		<-ctx.Done()

		// даем ему 30 секунд на остановку
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// наша основная задача, дождаться когда завершится эта функция, и лишь
		// потом выходить из функции

		// todo
		//
		// а что если сервер так и не запустился??
		err := s.Shutdown(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		stop() // если сервер не удалось запустить запускаем процедуру остановки
		return fmt.Errorf("server stop: %w", err)
	}

	// теперь ожидаем остановки сервера
	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("server stop: server shutdown gouroutine: %w", err)
	}
	return nil
}
