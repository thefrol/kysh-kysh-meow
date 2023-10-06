package main

import (
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
	"github.com/thefrol/kysh-kysh-meow/lib/scheduler"
)

func fileStorage(cfg config) (storage.Storager, error) {
	if cfg.FileStoragePath == "" {
		log.Info().Msg("Файл для сохранения и загрузки установлен в пустую строку, а значит все функции сохранения и загрузки на диск отключены")
		return storage.New(), nil
	}

	var s *storage.MemStore

	if cfg.Restore && fileExist(cfg.FileStoragePath) {
		var err error
		s, err = storage.FromFile(cfg.FileStoragePath)
		if err != nil {
			log.Error().Msgf("Не могу загрузить хранилища с диска %v по принчине %+v", cfg.FileStoragePath, err)
			return nil, err
		}
		log.Info().Msgf("Хранищиле воостановлено из %v", cfg.FileStoragePath)
	} else {
		ss := storage.New()
		s = &ss // todo помоему storage.New() должно возвращать указатель, чтобы не было таких проблем
	}

	// подключаем сохранение хранилища на диск

	wrapped := wrapStorageWithWrite(time.Duration(cfg.StoreIntervalSeconds)*time.Second, s, func() {
		s.ToFile(cfg.FileStoragePath)
		log.Info().Msg("Хранилище записано в файл")
	})
	return wrapped, nil
}

func fileExist(file string) bool {
	if s, err := os.Stat(file); err == nil && !s.IsDir() {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

// CallBackStorage обертывает хранилище так, что при изменении значений вызывает специальную коллбек функцию
// SaveCallBack, которую можно назначить. Можно использовать для синхронной записи на диск, при изменениях значений
type CallBackStorage struct {
	storage.Storager
	SaveCallback func()
}

func (s CallBackStorage) SetCounter(name string, value metrica.Counter) {
	s.Storager.SetCounter(name, value)
	s.SaveCallback()
}
func (s CallBackStorage) SetGauge(name string, value metrica.Gauge) {
	s.Storager.SetGauge(name, value)
	s.SaveCallback()
}

// wrapStorageWithWrite заботится о том, чтобы данные из хранилища сохранялись на диск. Есть два режима:
// синхронная запись(writeInterval==0), и накапливаемая запись( writeInterval>0). И для записи
// в том и том случае будет использовать функция callback
func wrapStorageWithWrite(writeInterval time.Duration, s storage.Storager, writeCallback func()) storage.Storager {
	// Если нужна синхронная запись, значит оборачиваем хранилище в CallBackStorage.
	// И при каждом сохранении счетчика записываем все на диск
	if writeInterval == 0 {
		cbs := CallBackStorage{Storager: s}
		cbs.SaveCallback = writeCallback
		log.Info().Msg("Создано хранилище с синхронной записью на диск")
		return cbs
	}

	if writeInterval < 500*time.Millisecond {
		log.Warn().Str("location", "server storage wrapper").Msgf("Указана слишком быстрое время сохранения метрик %vс. Это может сказать на производительности", writeInterval.Seconds())
	}
	// в ином случае запускаем планировщик
	sc := scheduler.New()
	sc.AddJob(writeInterval, writeCallback)
	go sc.Serve(200 * time.Millisecond) // todo это бы тоже как-то деликатно завершать надо, чтобы он случано остановился на половине записи при выключении
	return s
}
