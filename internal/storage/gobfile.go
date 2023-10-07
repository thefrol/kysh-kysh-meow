package storage

import (
	"encoding/gob"
	"os"

	"github.com/rs/zerolog/log"
)

// тут опять есть необходимость в пуле. Пул энкодеров джейсон, пул энкодеров gzip, я бы мог бы все это повторно использовать в разных местах программы

func (m MemStore) ToFile(fname string) error {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error().Msgf("Cant open file %v: %+v", fname, err)
		return err
	}
	err = gob.NewEncoder(file).Encode(&m)
	if err != nil {
		log.Error().Msgf("Cant marshal to gob %v: %+v", fname, err)
		return err
	}
	// вообще мы можем просто указывать маршалер и врайтер, и там че хочешь потом хоть джейсонь
	return nil
}

// RewriteFromFile перезаписывает хранилище  store данными из файла fname
func RewriteFromFile(fname string, store *MemStore) error {
	file, err := os.Open(fname)
	if err != nil {
		log.Error().Msgf("Cant open file %v: %+v", fname, err)
		return err
	}

	err = gob.NewDecoder(file).Decode(&store)
	if err != nil {
		log.Error().Msgf("Cant unmarshal from gob %v: %+v", fname, err)
		return err
	}

	return nil

	// TODO
	//
	// Варинаты. добавляет значения из файла, и может быть использует интерфейс Storager
	//
	// Самый главный вопрос, которым стоит руководствоваться: останутся ли мапы в стеке? Или декодер создаст новые памы которые сразу в кучу попадут?
	// Если так, то лучше сделать более долгую загрузку - пользоваться исходными мапами, просто переписать в них из исходного хранилища данные
}

// gob файл со всеми метриками 549 байт, без метрик 140
// json??? json+gzip?
