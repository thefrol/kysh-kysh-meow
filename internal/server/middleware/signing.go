package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/sign"
)

func Signing(key string) func(http.Handler) http.Handler {
	keyBytes := []byte(key)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedSign := r.Header.Get(sign.SignHeaderName)
			if receivedSign == "" {
				api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Нет подписи")
				return
			}

			if r.GetBody == nil {
				buf := bytes.NewBuffer(make([]byte, 0, 500))
				_, err := io.Copy(buf, r.Body)
				if err != nil {
					api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Cant replace request boby, signing failed")
					return
				}
				r.Body.Close()

				r.GetBody = func() (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
				}
				r.Body, _ = r.GetBody()
			}
			body, err := r.GetBody()
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Cant get request boby, signing check failed")
				return
			}
			defer body.Close()

			data := make([]byte, 500)
			n, err := body.Read(data)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "cant read body")
				return
			}

			if err := sign.Check(data[:n], keyBytes, receivedSign); err != nil {
				api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Подпись не прошла проверку")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
