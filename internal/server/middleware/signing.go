package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/sign"
	"github.com/thefrol/kysh-kysh-meow/lib/intercept"
)

const SignBufferSize = 1500

func Signing(key string) func(http.Handler) http.Handler {
	keyBytes := []byte(key)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			receivedSign := r.Header.Get(sign.SignHeaderName)
			if receivedSign == "" {
				//api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Нет подписи")
				//w.WriteHeader(http.StatusBadRequest)
				//w.Write([]byte("no sign"))
				//log.Info().Msg("Нет подписи")
				next.ServeHTTP(w, r)
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

			data := make([]byte, SignBufferSize)
			// bug если буфер меньше чем сообщение,
			// то типа не прочитается и подпись не сможет валидироваться
			// но большое буфер тоже не охота делать
			// TODO!!!
			n, err := body.Read(data)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "cant read body")
				return
			}

			if err := sign.Check(data[:n], keyBytes, receivedSign); err != nil {
				api.HTTPErrorWithLogging(w, http.StatusNotFound, "Подпись не прошла проверку")
				return
			}

			// теперь займемся ответом:
			buf := bytes.NewBuffer(data[:0])
			faker := intercept.WithBuffer(w, buf)

			next.ServeHTTP(faker, r)

			// теперь запишем все, что мы забуферизировали

			s, err := sign.Bytes(buf.Bytes(), keyBytes)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "не удается подписать ответ: %v", err)
				return
			}
			w.Header().Set(sign.SignHeaderName, s)
			log.Info().Str("sign", s).Msg("Запрос подписан")

			// записываем из буфера в оригинальный райтер
			if err := faker.Flush(); err != nil {
				api.HTTPErrorWithLogging(w, http.StatusNotFound, "Не переписать из фейк врайтера : %v", err)
				return
			}
		})
	}
}
