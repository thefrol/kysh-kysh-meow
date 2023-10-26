package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Bytes возвращает сжатый массив байт
func Bytes(data []byte, level int) ([]byte, error) {

	b := bytes.NewBuffer(make([]byte, 0, 500)) //todo нужна какая-то константа
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
