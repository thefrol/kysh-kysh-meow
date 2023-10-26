package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/sign"
	"github.com/thefrol/kysh-kysh-meow/lib/intercept"
)

func MeowLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			faker := intercept.WithBytesCounter(w)

			d := countTime(func() {
				// запустить обработку
				next.ServeHTTP(faker, r)
			})

			defer func() {
				msg := recover()
				if msg == nil {
					return
				}
				log.Error().
					Str("method", r.Method).
					Str("uri", r.RequestURI).
					Bool("gzippedRequest", encoded(r, "gzip")).
					Str("Sign", r.Header.Get(sign.SignHeaderName)). // todo sign.HeaderName
					Dur("Duration", d).
					Msgf("PANIC -> %v", msg)
			}()

			log.Info().
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Bool("gzippedRequest", encoded(r, "gzip")).
				Str("Sign", r.Header.Get(sign.SignHeaderName)). // todo sign.HeaderName
				Dur("Duration", d).
				Msg("Request ->")

			log.Info().
				Int("statusCode", faker.StatusCode()).
				Str("Content-Type", w.Header().Get("Content-Type")).
				Str("Sign", w.Header().Get(sign.SignHeaderName)).
				Int("Size", faker.BytesWritten()).
				// todo add gzipped response flag
				Msg("Response ->")
		})
	}
}

// countTime засекает время исполнения функции аргумента
func countTime(f func()) (d time.Duration) {
	defer func(t time.Time) {
		d = time.Since(t)
	}(time.Now())
	f()
	return
	// хрена се, эта функция работает примерно в тысячу раз быстрее чем прошлая с указателем
}
