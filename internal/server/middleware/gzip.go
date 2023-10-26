package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/lib/intercept"
)

var acceptedContentTypes = []string{
	"text/plain",
	"text/html",
	"text/css",
	"text/xml",
	"application/json",
	"application/javascript"}

const CompressionLevel = gzip.BestCompression

// GZIP это мидлварь для сервера, которая сжимает содержимое запроса
// если тело сообщения больше чем minLenж. Для сжатия первоначально
// создается буфер изначальной вместимости bufSize
func GZIP(minLen int, bufSize int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Eсли клиент поддерживает gzip, то подменяем врайтер, не забыв его закрыть, и отправляем
			// запрос дальше по цепочке
			if !acceptsEncoding(r, "gzip") {
				//не обжимаем
				next.ServeHTTP(w, r)
				return
			}

			// буферизируем
			buf := bytes.NewBuffer(make([]byte, 0, bufSize)) // todo что если у нас есть пул буферов?
			faker := intercept.WithBuffer(w, buf)
			next.ServeHTTP(faker, r)

			notZippable := buf.Len() < minLen ||
				faker.StatusCode() >= 300 ||
				!contentTypeZippable(w.Header().Get("Content-Type"))

			if notZippable {
				// записываем все в оригинальный врайтер, не сжимая
				faker.Flush()
				return
			}

			// мы решили ужать ответ,
			// для начала запишем новый заголовок,
			// и запишем код ответа

			w.Header().Add("Content-Encoding", "gzip")
			w.Header().Del("Content-Length") // очень важно, иначе будет ошибка при попытке закрыть врайтер компрессора

			if faker.StatusCode() != 0 {
				w.WriteHeader(faker.StatusCode())
			}

			gz, err := gzip.NewWriterLevel(w, CompressionLevel)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "GZIP writer init failed %v", err)
				return
			}

			n, err := io.Copy(gz, buf)
			fmt.Println(n)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Compressing %v", err) // todo нужны хелперы вроду api.InternalError
				return
			}

			if buf.Len() > 0 {
				log.Error().
					Msg("В буфере остались непрочитанные байты") // мой страх почему-то
			}

			err = gz.Close()
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Не могу записать в зиппер: %v", err)
			}

		})
	}
}

// acceptsEncoding возвращает true, если на основании запроса r
// можно сказать, клиент поддерживает кодирование в encoding
//
// if acceptsEncoding(r, "gzip"){...}
func acceptsEncoding(r *http.Request, encoding string) bool {
	for _, v := range r.Header.Values("Accept-Encoding") {
		if strings.Contains(v, encoding) {
			return true
		}
	}
	return false
}

// contentTypeZippable проверяет подходит ли Content-Type
// для сжатия
func contentTypeZippable(s string) bool {
	for _, ct := range acceptedContentTypes {
		if strings.Contains(s, ct) {
			return true
		}
	}
	return false
}

// todo
//
// для части вспомогательных функций и констант, я бы
// воспользовался internal/compress, и надо бы придумать
// семантику вызова:
// compress.CheckContentType? compess.Recommended?
// compress.Gzip() compress.Unzip() ?
// archive.Compress archive. Decompress?
// archive.RequestEncoded?
// как много вообще этот пакет должен знать о http???
