package mem

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrorNilSaver = fmt.Errorf("обращение к нулевому сейверу: %w", ErrorNilRef)
)

type IntervalicSaver struct {
	File     string
	Interval time.Duration
	Store    *MemStore

	started bool
	stop    chan struct{} // закроется, когда пора остановиться, стремно, наверное это должна быть функция
	stopped chan struct{} // закроектся, когда остановимся

	Log zerolog.Logger
}

// Run запускает процесс записи в файл, каждые Interval секунд
// в файл будут сброшены все метрики.
//
// Если err!=nil значит функция запустилась и должна быть
// остановлена функцией Stop(). Если же err!=nil, то
// интервальная запись не была запущена. И можно не оставливать
func (is *IntervalicSaver) Run() error {
	if is == nil {
		return fmt.Errorf("запуск записывателя: %w", ErrorNilSaver)
	}

	if is.Store == nil {
		return fmt.Errorf("запуск записывателя: %w", ErrorBadConfig)
	}

	if int64(is.Interval) <= 0 {
		return fmt.Errorf("интервал должен быть больше нуля: %w", ErrorBadConfig)
	}

	if is.File == "" {
		return fmt.Errorf("нужно указать непустой файл: %w", ErrorBadConfig)
	}

	// настроим каналы для начала и остановки
	is.stop = make(chan struct{})
	is.stopped = make(chan struct{})

	// чтобы у нас горутина остановки случано не запивла с пустого канала
	// укажем, что хранилище было запущено
	is.started = true

	// по этому тикеру мы будем следить
	// когда сохраняться
	t := time.NewTicker(is.Interval)

	go func() {
		// как только эта горутина закончится,
		// мы просигнализируем, что все остновилось
		defer close(is.stopped)

		// Основной цикл
	loop:
		for {
			select {

			// если сработал таймер - записываем
			case <-t.C:
				err := is.Store.Dump(is.File)
				if err != nil {
					is.Log.Error().
						Err(err).
						Str("file", is.File).
						Msg("не могу сохранить файл")
				}

			// иначе завершаем
			case <-is.stop:
				err := is.Store.Dump(is.File) // сохранимся обязательно ещё раз перед выходом
				if err != nil {
					is.Log.Error().
						Err(err).
						Str("file", is.File).
						Msg("не могу сохранить файл")
				}
				break loop
			}
		}
		log.Info().Msg("Хранилище остановлено")

	}()

	return nil
}

// Stop останавливает интервальную запись и дожидается
// произведения последней записи. В случае успеха вернётся
// err!=nil. Если интервальная запись не была запущена,
// вернет ErrorNotStarted
//
// todo было бы ещё неплохо указать дедлайт остановки через контекст
func (is *IntervalicSaver) Stop() error {
	if is == nil {
		return fmt.Errorf("остановка записывателя: %w", ErrorNilSaver)
	}

	if !is.started {
		is.Log.Error().
			Err(ErrorNotStarted).
			Msg("остановка записывателя")
		return fmt.Errorf("остановка записывателя: %w", ErrorNotStarted)
	}

	// сообщаем горутине остановиться
	log.Info().Msg("IntervalicSaver останавливается")
	close(is.stop)

	// ждем, когда операции будут закончены
	<-is.stopped
	log.Info().Msg("IntervalicSaver остановлен")

	return nil
}
