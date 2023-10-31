package graceful

import (
	"context"
	"fmt"
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
	server := http.Server{Addr: addr, Handler: router}

	// запустим горутину, которая будет слушать сигналы от системы, и при получении
	// начнет процедуру остановки сервера
	go func() {
		<-RequestStop()
		log.Debug().Msg("server wants to shut down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ShutdownGracefullyContext(shutdownCtx, &server)
	}()

	log.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Msgf("^0^ не могу запустить сервер: %v \n", err)
	}
	log.Error().Msg("Run() остановлен")

	// если ошибка при запуске сервера, то горутина не не получит сигнал, но в общем её вырубит система как бэ
}

// RequestStop возвращает канал через который придёт сообщение, что операционная система запросила завершение работы
func RequestStop() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return sig
}

// WithSignal возвращает контекст, который будет остановлен
// по запросу операционной системы, для таких сигналов как
// syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT
func WithSignal(ctx context.Context) context.Context {
	stoppable, cancel := context.WithCancel(ctx)
	sig := RequestStop()

	go func() {
		<-sig
		cancel()
		fmt.Println("горутина, созданная WithSignal(), остановлена")
	}()

	return stoppable
}

func ShutdownGracefullyContext(ctx context.Context, serv *http.Server) {
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatal().Msg("Время вышло, сервер будет завершен принудительно")
			return
		}
	}()

	// Trigger graceful shutdown
	err := serv.Shutdown(ctx)
	if err != nil {
		log.Fatal().Msg(err.Error())
		panic(err)
	}
	log.Info().Msg("^-^ рутина остановки сервера завершилась")
}
