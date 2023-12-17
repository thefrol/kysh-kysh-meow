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

// Это размер буфера, который будет использован
// для перехвата тела ответа. Чтобы лишний раз не
// гонять память, можно указать какой-то размер буфера
// который не надо будет аллоцировать.
//
// В идеале сюда должен полностью поместиться стандартный
// ответ сервера
const MinBufferSize = 15

// CheckSignature это мидлварь, которая проверяет подписи полученных
// запросов, используя пакет sign и ключ шифрования key. Подпись
// мы получаем из заголовка sign.SigningHeaderName
func CheckSignature(key string) func(http.Handler) http.Handler {
	keyBytes := []byte(key)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// получаем подпись из заголовков
			receivedSign := r.Header.Get(sign.SignHeaderName)

			// так получилось, в тестах, если запрос
			// передал пустой ключ подписи, то мы
			// не проверяем подпись, даже если у сервера
			// такой ключ указан
			if receivedSign == "" {
				next.ServeHTTP(w, r)
				return
			}

			// подменим тело запроса
			// то есть прочитаем все из тела запроса
			// потом запишем это все обратно
			body, err := io.ReadAll(r.Body)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Не удалось прочитать тело запроса")
				return
			}

			// закрываем тело исходного запроса
			err = r.Body.Close()
			if err != nil {
				log.Error().Msg("Не могу закрыть тело сообщения в подписывающей мидвари")
				return
			}

			// и подменяем тело исходного запроса
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// проверим подпись запроса
			if err := sign.Check(body, keyBytes, receivedSign); err != nil {
				api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Подпись не прошла проверку")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SignResponse подписывает ответы сервера ключом,
// получившуюся подпись будет положена в заголовок
// sign.SignHeaderName
func SignResponse(key string) func(http.Handler) http.Handler {
	keyBytes := []byte(key)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Перехватим все, что последующие обработчики запишут
			// в ответ. И это все будет лежать в buf
			buf := bytes.NewBuffer(make([]byte, 0, MinBufferSize))
			faker := intercept.WithBuffer(w, buf)

			// запускаем дальше по цепочке обработчики
			// результат их работы будет записан в faker
			next.ServeHTTP(faker, r)

			// в buf хранится буферизированный ответ,
			// теперь мы посчитаем подпись для него
			s, err := sign.Bytes(buf.Bytes(), keyBytes)
			if err != nil {
				api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "не удается подписать ответ: %v", err)
				return
			}

			// Теперь запишем в заголовки ответа подпись
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
