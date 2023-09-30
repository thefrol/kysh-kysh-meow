package storage

import (
	"encoding/gob"
	"os"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

// тут опять есть необходимость в пуле. Пул энкодеров джейсон, пул энкодеров gzip, я бы мог бы все это повторно использовать в разных местах программы

func (s MemStore) ToFile(fname string) error {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		ololog.Error().Msgf("Cant open file %v: %+v", fname, err)
		return err
	}
	err = gob.NewEncoder(file).Encode(&s)
	if err != nil {
		ololog.Error().Msgf("Cant marshal to gob %v: %+v", fname, err)
		return err
	}
	// вообще мы можем просто указывать маршалер и врайтер, и там че хочешь потом хоть джейсонь
	return nil
}

func FromFile(fname string) (*MemStore, error) {
	file, err := os.Open(fname)
	if err != nil {
		ololog.Error().Msgf("Cant open file %v: %+v", fname, err)
		return nil, err
	}

	s := New()
	err = gob.NewDecoder(file).Decode(&s)
	if err != nil {
		ololog.Error().Msgf("Cant unmarshal from gob %v: %+v", fname, err)
		return nil, err
	}
	// вообще мы можем просто указывать маршалер и врайтер, и там че хочешь потом хоть джейсонь
	return &s, nil
}

// gob файл со всеми метриками 549 байт, без метрик 140
// json??? json+gzip?
