package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

var ErrorRestoreFileNotExist = fmt.Errorf("файла, чтобы выгрузить хранилише не существует")

// FileStorage позволяет писать и восстанавливаться из файда
// при помощи функций Dump() и Restore(). Является оберткой над
// типом memStore
type FileStorage struct {
	MemStore
	FileName string
}

// NewFileStorage создает FileStorage из MemStore, таким образом
// Позволяя импользовать функции записывать и читать хранилища из файла
// fileName при помощи фунцкий Dump() и Restore()
func NewFileStorage(m *MemStore, fileName string) FileStorage {
	return FileStorage{MemStore: *m, FileName: fileName}
}

// Restore загружает в хранилища данные из FileName, при этом
// тукущие значения стираются
func (s FileStorage) Restore() error {
	if !fileExist(s.FileName) {
		return ErrorRestoreFileNotExist
	}
	err := RewriteFromFile(s.FileName, &s.MemStore)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла %v при загрузке хранилища: %v", s.FileName, err)
	}
	return nil
	// todo
	//
	// прям сюда весь код из MemStore по сохранению и загрузке
}

func (s FileStorage) Dump() error {
	return s.MemStore.ToFile(s.FileName)
}

// fileExist проверяет существует ли файл file, если да
// то возвращает true. Так же проверяет, что file не является
// директорией.
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

// IntervalDump это хранилище метрик, которое сохраняется в файл с заданной переодичностью Interval.
// Требудет запуска при помощи StartDumping()
type IntervalDump struct {
	FileStorage
	Interval time.Duration
}

// NewIntervalDump создает хранилище с записью через равные промеждутки времени interval,
// оборачивает уже существующее файловое хранилище s типа FileStorage.
func NewIntervalDump(s *FileStorage, interval time.Duration) IntervalDump {
	return IntervalDump{FileStorage: *s, Interval: interval}
	// TODO
	//
	// Стоит ли MemStore принимать в аргументе, поможет ли это держать все в стеке?
}

// StartDumping начинает процесс сохранения в файл через промежутки времени, и занимает
// текущий поток. Будет ждать отмены через ctx.Context.
func (s IntervalDump) StartDumping(ctx context.Context) {
	t := time.NewTicker(s.Interval)
	defer t.Stop()
	// todo а как закрывать канал?

	for {
		select {
		case <-t.C:
			// когда тикает таймер, мы просто записываем данные в хранилище. Если не получилось - ничего страшного
			err := s.Dump()
			if err != nil {
				log.Error().Msgf("Не удалось записать хранилише в файл %v: %+v", s.FileName, err)
				continue
			}
			log.Info().Msgf("Хранилище записано в %v", s.FileName)
		case <-ctx.Done():
			// Если пришел запрос на отключение - тоже записываем в файл
			err := s.Dump()
			if err != nil {
				log.Error().Msgf("При завершении работы IntervalDumper, не удалось записать хранилише в файл %v: %+v", s.FileName, err)
			}
			log.Info().Msgf("Запись по интервалу прекращает свою работу")
			return
		}
	}
}

// SyncDump - это хранилище с синхронной записью на диск, при изменении хоть одной метрики,
// все хранилище будет экспортировано.
type SyncDump struct {
	FileStorage
}

// NewSyncDump создает хранилище с сихнронной записью из
// экспортируемого хранилища s
func NewSyncDump(s *FileStorage) SyncDump {
	return SyncDump{FileStorage: *s}
}

// TODO
//
// Интерфейс может быть ExporterStorager, нужны только Dump()
//
// Может быть объединить FileStorage и SyncStorage?
//
// Может тут тоже сделать отдельный поток, пусть он пишет по сигналу, но не занимая время сервера?

// SetCounter устанавливает значение счетчика и инициализирует запись на диск
func (s SyncDump) SetCounter(name string, value metrica.Counter) {
	s.FileStorage.SetCounter(name, value)
	err := s.Dump()
	if err != nil {
		log.Error().Msgf("Не удалось записать в хранилише %v: %v", s.FileName, err)
	}
	log.Info().Msgf("После обновления метрики %v[Counter], хранилище записино в %v", name, s.FileName)

}

// SetGauge устанавливает значение счетчика типа gauge и инициализирует запись на диск
func (s SyncDump) SetGauge(name string, value metrica.Gauge) {
	s.FileStorage.SetGauge(name, value)
	err := s.Dump()
	if err != nil {
		log.Error().Msgf("Не удалось записать в хранилише %v: %v", s.FileName, err)
	}
	log.Info().Msgf("После обновления метрики %v[Gauge], хранилище записино в %v", name, s.FileName)
}
