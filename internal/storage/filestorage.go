package storage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

var ErrorRestoreFileNotExist = fmt.Errorf("файла, чтобы выгрузить хранилише, не существует")

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
	file, err := os.Open(s.FileName)
	if err != nil {
		log.Error().Msgf("Cant open file %v: %+v", s.FileName, err)
		return err
	}

	err = gob.NewDecoder(file).Decode(&s.MemStore)
	if err != nil {
		log.Error().Msgf("Cant unmarshal from gob %v: %+v", s.FileName, err)
		return err
	}

	return nil

	// TODO
	//
	// Вариант: добавляет значения из файла, не очищая хранилище, и может быть использует интерфейс Storager
	//
	// Самый главный вопрос, которым стоит руководствоваться: останутся ли мапы в стеке? Или декодер создаст новые памы которые сразу в кучу попадут?
	// Если так, то лучше сделать более долгую загрузку - пользоваться исходными мапами, просто переписать в них из исходного хранилища данные

}

func (s FileStorage) Dump() error {
	file, err := os.OpenFile(s.FileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error().Msgf("Cant open file %v: %+v", s.FileName, err)
		return err
	}
	err = gob.NewEncoder(file).Encode(&s.MemStore)
	if err != nil {
		log.Error().Msgf("Cant marshal to gob %v: %+v", s.FileName, err)
		return err
	}
	// вообще мы можем просто указывать маршалер и врайтер, и там че хочешь потом хоть джейсонь
	// например, могут быть функопции с настройками декодеров и энкодеров
	return nil
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
