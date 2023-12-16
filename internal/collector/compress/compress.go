package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/rs/zerolog/log"
)

// todo по версии влада избавиться от этой темы с длинной
// вообще эту темку можно и убрать
const Bufferlen = 500

// Bytes возвращает сжатый массив байт
func Bytes(data []byte, level int) ([]byte, error) {

	b := bytes.NewBuffer(make([]byte, 0, Bufferlen))
	gz, err := gzip.NewWriterLevel(b, level)
	if err != nil {
		return nil, fmt.Errorf("cant create compressor")
	}

	_, err = gz.Write(data)
	if err != nil {
		return nil, fmt.Errorf("cant write to zip")
	}

	err = gz.Close()
	if err != nil {
		return nil, fmt.Errorf("ошибка компрессии и закрытия компрессора")
	}

	log.Info().
		Int("size_before", len(data)).
		Int("size_after", b.Len()).
		Float64("compression_ratio", float64(b.Len())/float64(len(data))).
		Msg("Компрессор закончил работать")

	return b.Bytes(), nil
}

// todo
//
// Думаю, именно в этом пакете я бы хотел объявить пулы с gzip.Ecoder и gzip.decoder
