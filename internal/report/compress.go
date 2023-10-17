package report

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const compressMinLenght = 20

func init() {
	AddMiddleware(ApplyGZIP(compressMinLenght, gzip.BestCompression))
}

func ApplyGZIP(minLenght int, level int) func(c *resty.Client, r *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		// проверяем, что контент уже не закодирован каким-нибудт другим мидлварью или кодом
		if r.Header.Values("Content-Encoding") != nil {
			log.Warn().Str("location", "agent/middleware/gzip").Fields(r.Header).Msg("Запрос на сервер уже сжат")
			return nil
		}

		v, ok := r.Body.(io.Reader)
		if !ok {
			log.Warn().Str("location", "agent/middleware/gzip").Msg("Тело сообщение передано не в формате io.Reader, а значит мы не можем его заархивировать. Вообще мы хотим передавать тело сообщения именно ридером")
			return nil
		}

		// теперь мы прочитаем minLenght байт из буфера
		// Если тело к этому моменту закончится, то отправляем без сжатия
		// А если нет, то пишем в gzip уже прочитанное, и все остальное

		t := make([]byte, minLenght)
		n, _ := io.ReadAtLeast(v, t, minLenght)
		if n < minLenght {
			// не забудем обрезать буфер. Мы возьмем оттуда
			// только первые n символов, потому что остальные будут просто нулями
			// вплоть до minLenght
			r.SetBody(bytes.NewBuffer(t[:n]))
			log.Info().Msg("Request too short for compressing, sending as is")
			return nil
		}

		// тут мы часть буфера прочитали, и хотим прочитать оставшееся, и скомпрессировать все вместе

		b := new(bytes.Buffer) // todo мы можем писать сразу в тело сообщения я думаю при желании, возможноо.... И если не будем выводить инфу по сжатию
		gz, err := gzip.NewWriterLevel(b, level)
		if err != nil {
			return fmt.Errorf("cant create compressor")
		}

		bytesBefore := n // посчитаем успешность нашей компрессии, для этого запомним сколько байт было изначально

		_, err = gz.Write(t)
		if err != nil {
			return fmt.Errorf("cant write min lenght buffer back to zipper")
		}

		nAfter, err := io.Copy(gz, v)
		if err != nil {
			return fmt.Errorf("cant write to gzip writer")
		}
		bytesBefore += int(nAfter)

		r.Header.Add("Content-Encoding", "gzip")
		r.SetBody(b)
		gz.Close()

		log.Info().
			Int("size_before", bytesBefore).
			Int("size_after", b.Len()).
			Float64("compression_ratio", float64(b.Len())/float64(bytesBefore)).
			Msg("Компрессор закончил работать")

		return nil
	}
}
