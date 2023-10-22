package report

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const compressMinLenght = 20

func init() {
	AddMiddleware(ApplyGZIP(compressMinLenght, gzip.BestCompression))
}

// Вообще рести не очень подходит под клиента тут, если я собираюсь использовать
// пакет retry вместо resty.WithRetry(), потому что поведение может
// быть довольно неожиданным.
//
// Например, если мы мидвари запускаем перед каждый запуском, а тело сообщения
// уже сформировано, то он может подписать сначала данные чистые, потом
// подпишет сжатые.
//
// Мне бы хотелось иметь такой пакет, где запрос составляется один раз, и далее
// без изменений отправляется. Или я должен тут мидлвари переписывать, или
// отказываться от мидлварей вообще, что тоже вариант, конечно. Или свой фреймворк
// мидлварей писать, чего я не хочу

func ApplyGZIP(minLenght int, level int) func(c *resty.Client, r *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		// проверяем, что контент уже не закодирован каким-нибудт другим мидлварью или кодом
		if r.Header.Values("Content-Encoding") != nil {
			log.Warn().Str("location", "agent/middleware/gzip").Fields(r.Header).Msg("Запрос на сервер уже сжат")
			return nil
		}

		v, ok := r.Body.([]byte)
		if !ok {
			log.Warn().Str("location", "agent/middleware/gzip").Msg("Тело сообщение передано не в формате []byte, а значит мы не можем его заархивировать. Вообще мы хотим передавать тело сообщения именно массивом байт")
			return nil
		}

		// посмотрим, если нам стоит сжимать сообщение, проверим его длинну
		if len(v) < minLenght {
			log.Info().Msg("Request too short for compressing, sending as is")
			return nil
		}

		// тут мы часть буфера прочитали, и хотим прочитать оставшееся, и скомпрессировать все вместе

		b := bytes.NewBuffer(make([]byte, 0, 500)) //todo нужна какая-то константа
		gz, err := gzip.NewWriterLevel(b, level)
		if err != nil {
			return fmt.Errorf("cant create compressor")
		}

		_, err = gz.Write(v)
		if err != nil {
			return fmt.Errorf("cant write to zip")
		}

		gz.Close()

		r.Header.Add("Content-Encoding", "gzip")
		r.SetBody(b.Bytes())
		gz.Close()

		log.Info().
			Int("size_before", len(v)).
			Int("size_after", b.Len()).
			Float64("compression_ratio", float64(b.Len())/float64(len(v))).
			Msg("Компрессор закончил работать")

		return nil
	}
}
