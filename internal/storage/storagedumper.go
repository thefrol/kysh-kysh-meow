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

func NewFileStorage(m *MemStore, fileName string) FileStorage {
	return FileStorage{MemStore: *m, FileName: fileName}
}

// Перезаписать данные значением из файла fname
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

func (s FileStorage) Dump() error {
	return s.MemStore.ToFile(s.FileName)
}

type IntervalDump struct {
	FileStorage
	Interval time.Duration
}

func NewIntervalDump(s *FileStorage, interval time.Duration) IntervalDump {
	return IntervalDump{FileStorage: *s, Interval: interval}
	// TODO
	//
	// Стоит ли MemStore принимать в аргументе, поможет ли это держать все в стеке?
}

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

type SyncDump struct {
	FileStorage
}

func NewSyncDump(s *FileStorage) SyncDump {
	return SyncDump{FileStorage: *s}
}

func (s SyncDump) SetCounter(name string, value metrica.Counter) {
	s.FileStorage.SetCounter(name, value)
	err := s.Dump()
	if err != nil {
		log.Error().Msgf("Не удалось записать в хранилише %v: %v", s.FileName, err)
	}
	log.Info().Msgf("После обновления метрики %v[Counter], хранилище записино в %v", name, s.FileName)

}
func (s SyncDump) SetGauge(name string, value metrica.Gauge) {
	s.FileStorage.SetGauge(name, value)
	err := s.Dump()
	if err != nil {
		log.Error().Msgf("Не удалось записать в хранилише %v: %v", s.FileName, err)
	}
	log.Info().Msgf("После обновления метрики %v[Gauge], хранилище записино в %v", name, s.FileName)
}
