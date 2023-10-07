// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

func main() {
	cfg := mustConfigure(defaultConfig)

	// создаем хранилище
	m := storage.New()
	s, cancelStorage := ConfigureStorage(&m, cfg)

	// Запускаем сервер с поддержкой нежного завершения
	Run(cfg, s)

	// Завершаем последние дела
	// попытаемся сохраниться в файл
	cancelStorage()

	// Даем ему время
	time.Sleep(time.Second)

	log.Info().Msg("^.^ Сервер завершен нежно")
	// Wait for server context to be stopped

}

func Run(cfg config, s storage.Storager) {
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
	}

	<-serverCtx.Done()
}

func ConfigureStorage(m *storage.MemStore, cfg config) (storage.Storager, context.CancelFunc) {
	// 1. Если путь не задан, то возвращаем хранилище в оперативке, без приблуд
	// 2. Иначе оборачиваем файловым хранилищем, но не возвращаем пока
	// 3. Если Restore=true, то читаем из файла. Если файла не существует, то игнорируем проблему
	// 4. Оборачиваем файловое хранилище сихнронным или интервальным сохранением

	// Если путь до хранилища не пустой, то нам нужно инициаизировать обертки над хранилищем
	if cfg.FileStoragePath == "" {
		return m, func() {
			log.Info().Msg("Хранилище сихронной записи получило сигнал о завершении, но файловая запись в текущей конфигурации сервера не используется. Ничего не записано")
		}
	}

	// Оборачиваем файловым хранилищем, в случае, есл и
	fs := storage.NewFileStorage(m, cfg.FileStoragePath)
	if cfg.Restore {
		err := fs.Restore()
		// в случае, если файла не существует, игнорируем эту проблему
		if err != nil && err != storage.ErrorRestoreFileNotExist {
			panic(err)
		}
	}

	if cfg.StoreIntervalSeconds == 0 {
		// инициализируем сихнронную запись,
		// при этом сохраняться в конце нам не понадобится
		return storage.NewSyncDump(&fs), func() {
			log.Info().Msg("Хранилище сихронной записи получило сигнал о завершении, но дополнительно сохранение не нужно")
		}
	}

	// Запускаем интервальную запись и создаем токен отмены, при необходимости сюда можно будет добавить и группу ожидания
	s := storage.NewIntervalDump(&fs, time.Duration(cfg.StoreIntervalSeconds)*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go s.StartDumping(ctx)

	return s, func() {
		// оберстка сделана под группу ожидаения
		cancel()
	}

	// TODO
	//
	// Пока что судя по всему эта функция эвакуируем моё хранилище из стека, мне кажется, но не мапы
	// Как вариант сразу создавать FileStorage, и просто не оборачивать его если надо
}
