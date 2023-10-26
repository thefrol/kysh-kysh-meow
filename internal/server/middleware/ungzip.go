package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// UnGZIP распаковывает запросы, закодированные при помощи GZIP, и пропускает все
// остальное мимо ушей, отслеживает чтобы хотя бы один из заголовков Content-Encoding
// содержал подстроку gzip
func UnGZIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// пропускаем, поскольку мы не обрабатываем такие сжималки
		if !encoded(r, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// в случае если перед нами gzip, заменяем исходное тело запроса на обертку с gzip
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Не могу декомпрессировать тело запроса %v", err)
			return
		}
		defer r.Body.Close()
		r.Body = gz //todo: make a pool of zgips and return only buffer, not gzipeers

		next.ServeHTTP(w, r)

	})
}

// encoded возвращает true, если запрос r закодирован
// кодировщиком encoder. Например:
//
// if encoded(r,"gzip"){...}
func encoded(r *http.Request, encoder string) bool {
	// По стандартным договоренностям, если запрос сжат или закодирован, то кодировщики
	// указываются в последовательности их применения, а значит нам нужно читать последний
	// заголовок Content-Encoded, и если он gzip, то расшифровать

	hh := r.Header.Values("Content-Encoding")
	// Если вообще нет таких заголовков, то возвращаем false
	if len(hh) == 0 {
		return false
	}

	lastValue := hh[len(hh)-1]
	return strings.Contains(lastValue, encoder)
}
