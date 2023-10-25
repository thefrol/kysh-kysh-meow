// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	log.Info().Msgf("Сервер запущен строкой %v", strings.Join(os.Args, " "))

	cfg := mustConfigure(defaultConfig)

	// создаем хранилище
	s, cancelStorage := ConfigureStorage(cfg)

	// Создаем объект App, который в дальнейшем включит в себя все остальное тут
	// app, err := app.New(context.TODO(), cfg.DatabaseDSN)
	// if err != nil {
	// 	log.Fatal().Msgf("Ошибка во время конфигурирования сервера %v", err)
	// 	panic(err)
	// }
	// if err := app.CheckConnection(context.Background()); err == nil {
	// 	log.Info().Msg("Связь с базой данных в порядке")
	// }

	// Запускаем сервер с поддержкой нежного завершения,
	// занимаем текущий поток до вызова сигнатов выключения
	Run(cfg, s)

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
func Run(cfg config, s api.Operator) {
	// Запускаем сервер с поддержкой нежного выключения
	// вдохноввлено примерами роутера chi
	server := http.Server{Addr: cfg.Addr, Handler: router.MeowRouter(s, cfg.Key)}

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
		// если не биндится, то хотя бы выходить с ошибкой
		//
		// можно дать несколько попыток забиндиться
	}

	<-serverCtx.Done()

	// МНе кажется в отдельную функцию надо выделить именно все, что относится к нежному завершению, + надо перевести комменты по коду на русский
}

// ConfigureStorage подготавливает хранилище к работе в соответствии с текущими настройками,
// при необходимости загружает из файла значения метрик, запускает сохранение в файл, и
// возвращает интерфейс хранилища и функцию, подготавливающая ханилище к остановке
//
// На входе получает экземпляр хранилища m, и далее оборачивает его другим классов,
// наиболее соответсвующим задаче, исходя из cfg
func ConfigureStorage(cfg config) (api.Operator, context.CancelFunc) {
	// 0. Если указана база данных, создаем хранилище с базой данных
	// 1. Если путь не задан, то возвращаем хранилище в оперативке, без приблуд
	// 2. Иначе оборачиваем файловым хранилищем, но не возвращаем пока
	// 3. Если Restore=true, то читаем из файла. Если файла не существует, то игнорируем проблему
	// 4. Оборачиваем файловое хранилище сихнронным или интервальным

	// TODO
	//
	// Думаю эта функция должна получать на вход контекст, а не возвращать CancelFunc
	//
	// И нужно убрать отсюда все паники

	// Если база данных
	if cfg.DatabaseDSN != "" {
		db, err := sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			e := fmt.Errorf("не могу создать соединение с БД: %w", err)
			panic(e)
		}

		dbs, err := storage.NewDatabase(db)
		if err != nil {
			log.Error().Msgf("Ошибка создания хранилища в базе данных - %v", err)
			panic(err)
		}

		if err := dbs.Check(context.TODO()); err != nil {
			log.Error().Msgf("Нет соединения с БД - %v", err)
		}

		log.Info().Msg("Создано хранилише в Базе данных")
		return dbs, func() {
			err := db.Close()
			if err != nil {
				log.Error().Msgf("Не могу закрыть базу данных: %v", err)
			}

			// todo
			//
			// Конечно, я хочу делать это defer или как-то так, можно у нас будет некий app.Close()
		}
	}

	// Если не база данных, то начинаем с начала - создаем хранилище в памяти, и оборачиваем его всякими штучками если надо
	m := storage.New()

	// Если путь до хранилища не пустой, то нам нужно инициаизировать обертки над хранилищем
	if cfg.FileStoragePath == "" {
		log.Info().Msg("Установлено хранилище в памяти. Сохранение на диск отключено")
		return storage.AsOperator(m), func() {
			log.Info().Msg("Хранилище сихронной записи получило сигнал о завершении, но файловая запись в текущей конфигурации сервера не используется. Ничего не записано")
		}
	}

	// Оборачиваем файловым хранилищем, в случае, есл и
	fs := storage.NewFileStorage(&m, cfg.FileStoragePath)
	if cfg.Restore {
		err := fs.Restore()
		// в случае, если файла не существует, игнорируем эту проблему
		if err != nil && err != storage.ErrorRestoreFileNotExist {
			panic(err)
		}
		log.Info().Msgf("Значения метрик загружены из %v", fs.FileName)
	}

	if cfg.StoreIntervalSeconds == 0 {
		// инициализируем сихнронную запись,
		// при этом сохраняться в конце нам не понадобится
		log.Info().Msgf("Установлено синхронное сохранение в %v в при записи", fs.FileName)
		return storage.AsOperator(storage.NewSyncDump(&fs)), func() {
			log.Info().Msg("Хранилище сихронной записи получило сигнал о завершении, но дополнительно сохранение не нужно")
		}
	}

	// Запускаем интервальную запись и создаем токен отмены, при необходимости сюда можно будет добавить и группу ожидания
	s := storage.NewIntervalDump(&fs, time.Duration(cfg.StoreIntervalSeconds)*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go s.StartDumping(ctx)

	log.Info().Msgf("Установлено сохранение с интервалом %v в %v в при записи", s.Interval, s.FileName)

	return storage.AsOperator(s), func() {
		// оберстка сделана под группу ожидаения
		cancel()
	}

	// TODO
	//
	// Пока что судя по всему эта функция эвакуируем моё хранилище из стека, мне кажется, но не мапы
	// Как вариант сразу создавать FileStorage, и просто не оборачивать его если надо
}
