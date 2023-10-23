package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/sign"
)

func Signing(key string) func(http.Handler) http.Handler {
	keyBytes := []byte(key)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json") // todo костылек убери лол

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

			data := make([]byte, 500)
			n, err := body.Read(data)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "cant read body")
				return
			}

			if err := sign.Check(data[:n], keyBytes, receivedSign); err != nil {
				log.Info().Str("receivedSign", receivedSign).Msg("подпись не прошла проверку")
				api.HTTPErrorWithLogging(w, http.StatusNotFound, "Подпись не прошла проверку")
				return
			}

			// теперь займемся запросом:
			// обернем в перехватчик врайтер и подпишем ответ
			fakew := SignInterceptor{
				WriteInterceptor: WriteInterceptor{
					origWriter: w,
					buf:        bytes.NewBuffer(data),
				},
				key: keyBytes,
				// todo переиспользуем наш буфер, а моем наверное целиком весь буфер если его почистить
			}

			next.ServeHTTP(&fakew, r)

			// todo  в целом мы могли бы вместо Close тут разобраться с перехваченными данными
			// Это было чуть более переиспользуемо и наверное понятно

		})
	}
}

type WriteInterceptor struct {
	statusCode int
	origWriter http.ResponseWriter
	buf        *bytes.Buffer
}

func (w *WriteInterceptor) WriteHeader(code int) {
	w.statusCode = code
}

func (w *WriteInterceptor) Header() http.Header {
	return w.origWriter.Header()
}

type SignInterceptor struct {
	WriteInterceptor
	key []byte
}

func (w WriteInterceptor) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

// Close говорит о том, что сообщение готовится к отправке, значит
// можно установить нужные хедеры и отправлять
func (w WriteInterceptor) Close() {
	if w.statusCode != 0 {
		w.origWriter.WriteHeader(w.statusCode)
	}

	_, err := io.Copy(w.origWriter, w.buf)
	if err != nil {
		log.Error().Msgf("copy a response to originalWriter: %v", err)
	}
}

func (w SignInterceptor) Close() {
	s, err := sign.Bytes(w.buf.Bytes(), w.key)
	if err != nil {
		log.Error().Msgf("cant sign response: %v", err)
	}

	// копируем все из временного хранилища по назначению
	w.WriteInterceptor.Header().Set(sign.SignHeaderName, s)
	log.Info().Str("sign", s).Msg("Запрос подписан")

	w.WriteInterceptor.Close()

}
