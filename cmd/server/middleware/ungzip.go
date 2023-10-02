package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
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
			ololog.Error().Str("location", "middleware/gzip").Strs("Content-Encoding", r.Header.Values("Content-Encoding")).Err(err).Msg("Cant unzip a request body")
			http.Error(w, "cant unzip", http.StatusBadRequest)
			//todo try recover and send data as is
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
	for _, v := range r.Header.Values("Content-Encoding") {
		if strings.Contains(v, encoder) {
			return true
		}
	}
	return false
}
